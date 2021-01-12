/*
Copyright 2021 The Kruise Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"context"
	"flag"
	"time"

	podutil "k8s.io/kubernetes/pkg/api/v1/pod"

	ctrlmeshproto "github.com/openkruise/kruise/apis/ctrlmesh/proto"
	ctrlmeshv1alpha1 "github.com/openkruise/kruise/apis/ctrlmesh/v1alpha1"
	ctrlmeshutil "github.com/openkruise/kruise/pkg/ctrlmesh/util"
	"github.com/openkruise/kruise/pkg/util"
	utildiscovery "github.com/openkruise/kruise/pkg/util/discovery"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	kubecontroller "k8s.io/kubernetes/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	concurrentReconciles = flag.Int("ctrlmesh-server-workers", 3, "Max concurrent workers for CtrlMesh Server controller.")
	controllerKind       = ctrlmeshv1alpha1.SchemeGroupVersion.WithKind("VirtualOperator")
)

func SetupWithManager(mgr manager.Manager) error {
	if !utildiscovery.DiscoverGVK(controllerKind) {
		return nil
	}
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) *ReconcileCtrlMeshServer {
	return &ReconcileCtrlMeshServer{
		Client: util.NewClientFromManager(mgr, "ctrlmesh-server-controller"),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r *ReconcileCtrlMeshServer) error {
	// Create a new controller
	c, err := controller.New("ctrlmesh-server-controller", mgr, controller.Options{Reconciler: r, MaxConcurrentReconciles: *concurrentReconciles})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &ctrlmeshv1alpha1.VirtualOperator{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1.Pod{}}, &podEventHandler{reader: mgr.GetCache()})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Channel{Source: grpcRecvTriggerChannel}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1.Namespace{}}, &namespaceEventHandler{reader: mgr.GetCache()})
	if err != nil {
		return err
	}

	return nil
}

type ReconcileCtrlMeshServer struct {
	client.Client
}

// +kubebuilder:rbac:groups=ctrlmesh.kruise.io,resources=virtualoperators,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=ctrlmesh.kruise.io,resources=virtualoperators/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch

func (r *ReconcileCtrlMeshServer) Reconcile(request reconcile.Request) (res reconcile.Result, retErr error) {
	startTime := time.Now()
	defer func() {
		if retErr == nil {
			if res.Requeue || res.RequeueAfter > 0 {
				klog.Infof("Finished syncing VOP %s, cost %v, result: %v", request, time.Since(startTime), res)
			} else {
				klog.Infof("Finished syncing VOP %s, cost %v", request, time.Since(startTime))
			}
		} else {
			klog.Errorf("Failed syncing VOP %s: %v", request, retErr)
		}
	}()

	vop := &ctrlmeshv1alpha1.VirtualOperator{}
	err := r.Get(context.TODO(), request.NamespacedName, vop)
	if err != nil {
		if errors.IsNotFound(err) {
			// TODO: need to reset all Pods?
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	pods, err := r.getPodsForVOP(vop)
	if err != nil {
		return reconcile.Result{}, err
	}
	for _, pod := range pods {
		if v, ok := expectationsSrvHash.Load(types.NamespacedName{Namespace: pod.Namespace, Name: pod.Name}); ok {
			klog.Warningf("Skip reconcile VOP %s for Pod %s has dirty expectation %s", request, pod.Name, util.DumpJSON(v))
			return reconcile.Result{}, nil
		}
	}
	namespaces, err := r.getActiveNamespaces()
	if err != nil {
		return reconcile.Result{}, err
	}

	protoRouteMap := generateProtoRoute(vop, namespaces)
	protoEndpoints := generateProtoEndpoints(vop, pods)
	for _, pod := range pods {
		var conn *grpcSrvConnection
		if v, ok := cachedGrpcSrvConnection.Load(types.NamespacedName{Namespace: pod.Namespace, Name: pod.Name}); !ok || v == nil {
			if podutil.IsPodReady(pod) {
				klog.V(4).Infof("VOP %s/%s find no connection from Pod %s yet", vop.Namespace, vop.Name, pod.Name)
			}
			continue
		} else {
			conn = v.(*grpcSrvConnection)
		}

		subset := determinePodSubset(vop, pod)

		newSpec := &ctrlmeshproto.Spec{
			Meta:      &ctrlmeshproto.SpecMeta{VopName: vop.Name, VopResourceVersion: vop.ResourceVersion},
			Route:     protoRouteMap[subset],
			Endpoints: protoEndpoints,
		}
		if err = r.handleSrvConnection(vop, pod, conn, newSpec); err != nil {
			klog.Errorf("Failed to handle srv connection for VOP %s/%s to Pod %s: %v", vop.Namespace, vop.Name, pod.Name, err)
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileCtrlMeshServer) handleSrvConnection(vop *ctrlmeshv1alpha1.VirtualOperator, pod *v1.Pod, conn *grpcSrvConnection, newSpec *ctrlmeshproto.Spec) error {
	conn.mu.Lock()
	prevStatus := conn.status
	conn.mu.Unlock()

	newSpecHash := ctrlmeshutil.CalculateHashForProtoSpec(newSpec)

	if prevStatus.specHash != nil &&
		prevStatus.specHash.RouteHash == newSpecHash.RouteHash &&
		prevStatus.specHash.EndpointsHash == newSpecHash.EndpointsHash &&
		prevStatus.specHash.NamespacesHash == newSpecHash.NamespacesHash {
		return nil
	}

	klog.Infof("Sending proto spec %v (hash: %v) to Pod %s in VOP %s/%s", util.DumpJSON(newSpec), util.DumpJSON(newSpecHash), pod.Name, vop.Namespace, vop.Name)
	expectationsSrvHash.Store(types.NamespacedName{Namespace: pod.Namespace, Name: pod.Name}, newSpecHash)
	if err := conn.srv.Send(newSpec); err != nil {
		expectationsSrvHash.Delete(types.NamespacedName{Namespace: pod.Namespace, Name: pod.Name})
		return err
	}
	return nil
}

func (r *ReconcileCtrlMeshServer) getPodsForVOP(vop *ctrlmeshv1alpha1.VirtualOperator) ([]*v1.Pod, error) {
	podList := v1.PodList{}
	if err := r.List(context.TODO(), &podList, client.InNamespace(vop.Namespace), client.MatchingLabels{ctrlmeshv1alpha1.VirtualOperatorInjectedKey: vop.Name}); err != nil {
		return nil, err
	}
	var pods []*v1.Pod
	for i := range podList.Items {
		pod := &podList.Items[i]
		if kubecontroller.IsPodActive(pod) {
			pods = append(pods, pod)
		}
	}
	return pods, nil
}

func (r *ReconcileCtrlMeshServer) getActiveNamespaces() ([]*v1.Namespace, error) {
	namespaceList := v1.NamespaceList{}
	if err := r.List(context.TODO(), &namespaceList); err != nil {
		return nil, err
	}
	var namespaces []*v1.Namespace
	for i := range namespaceList.Items {
		ns := &namespaceList.Items[i]
		if ns.DeletionTimestamp == nil {
			namespaces = append(namespaces, ns)
		}
	}
	return namespaces, nil
}
