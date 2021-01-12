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

package metainfo

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sort"
	"time"

	"github.com/openkruise/kruise/pkg/grpcregistry"

	publicv1alpha1 "github.com/openkruise/kruise/apis/public/v1alpha1"
	"github.com/openkruise/kruise/pkg/util"
	webhookutil "github.com/openkruise/kruise/pkg/webhook/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
	kubecontroller "k8s.io/kubernetes/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	podSelector labels.Selector
	namespace   string
	localName   string
)

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) *ReconcileMetaInfo {
	return &ReconcileMetaInfo{
		Client: util.NewClientFromManager(mgr, "metainfo-controller"),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r *ReconcileMetaInfo) error {
	namespace = webhookutil.GetNamespace()
	if localName = os.Getenv("POD_NAME"); len(localName) == 0 {
		return fmt.Errorf("find no POD_NAME in env")
	}

	// Read the service of kruise webhook, to get the pod selector
	svc := &v1.Service{}
	svcNamespacedName := types.NamespacedName{Namespace: namespace, Name: webhookutil.GetServiceName()}
	err := mgr.GetAPIReader().Get(context.TODO(), svcNamespacedName, svc)
	if err != nil {
		return fmt.Errorf("get service %s error: %v", svcNamespacedName, err)
	}

	podSelector, err = metav1.LabelSelectorAsSelector(&metav1.LabelSelector{MatchLabels: svc.Spec.Selector})
	if err != nil {
		return fmt.Errorf("parse service %s selector %v error: %v", svcNamespacedName, svc.Spec.Selector, err)
	}

	// Create a new controller
	c, err := controller.New("metainfo-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &publicv1alpha1.MetaInfo{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1.Pod{}}, &enqueueHandler{}, predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			pod := e.Object.(*v1.Pod)
			return podSelector.Matches(labels.Set(pod.Labels))
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			pod := e.ObjectNew.(*v1.Pod)
			return podSelector.Matches(labels.Set(pod.Labels))
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			pod := e.Object.(*v1.Pod)
			return podSelector.Matches(labels.Set(pod.Labels))
		},
		GenericFunc: func(e event.GenericEvent) bool {
			pod := e.Object.(*v1.Pod)
			return podSelector.Matches(labels.Set(pod.Labels))
		},
	})
	if err != nil {
		return err
	}

	return nil
}

type enqueueHandler struct{}

func (e *enqueueHandler) Create(_ event.CreateEvent, q workqueue.RateLimitingInterface) {
	q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
		Name: publicv1alpha1.MetaNameOfKruiseManager,
	}})
}

func (e *enqueueHandler) Update(_ event.UpdateEvent, q workqueue.RateLimitingInterface) {
	q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
		Name: publicv1alpha1.MetaNameOfKruiseManager,
	}})
}

func (e *enqueueHandler) Delete(_ event.DeleteEvent, q workqueue.RateLimitingInterface) {
	q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
		Name: publicv1alpha1.MetaNameOfKruiseManager,
	}})
}

func (e *enqueueHandler) Generic(_ event.GenericEvent, q workqueue.RateLimitingInterface) {
	q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
		Name: publicv1alpha1.MetaNameOfKruiseManager,
	}})
}

var _ reconcile.Reconciler = &ReconcileMetaInfo{}

// ReconcileMetaInfo reconciles a MetaInfo object
type ReconcileMetaInfo struct {
	client.Client
}

// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch
// +kubebuilder:rbac:groups=public.kruise.io,resources=metainfos,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=public.kruise.io,resources=metainfos/status,verbs=get;update;patch

func (r *ReconcileMetaInfo) Reconcile(request reconcile.Request) (res reconcile.Result, err error) {
	if request.Name != publicv1alpha1.MetaNameOfKruiseManager {
		klog.Infof("Ignore MetaInfo %s", request.Name)
		return reconcile.Result{}, nil
	}

	start := time.Now()
	klog.V(3).Infof("Starting to process MetaInfo %v", request.Name)
	defer func() {
		if err != nil {
			klog.Warningf("Failed to process MetaInfo %v, elapsedTime %v, error: %v", request.Name, time.Since(start), err)
		} else {
			klog.Infof("Finish to process MetaInfo %v, elapsedTime %v", request.Name, time.Since(start))
		}
	}()

	podList := &v1.PodList{}
	err = r.List(context.TODO(), podList, client.InNamespace(namespace), client.MatchingLabelsSelector{Selector: podSelector})
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("list pods in %s error: %v", namespace, err)
	}

	var hasLeader bool
	endpoints := make(publicv1alpha1.MetaInfoEndpoints, 0, len(podList.Items))
	for i := range podList.Items {
		pod := &podList.Items[i]
		if !kubecontroller.IsPodActive(pod) {
			continue
		}

		e := publicv1alpha1.MetaInfoEndpoint{Name: pod.Name, PodIP: pod.Status.PodIP}
		if pod.Name == localName {
			e.Leader = true
			hasLeader = true
		}
		endpoints = append(endpoints, e)
	}
	sort.Sort(endpoints)
	if !hasLeader {
		return reconcile.Result{}, fmt.Errorf("no leader %s in new endpoints %v", localName, util.DumpJSON(endpoints))
	}

	ports := publicv1alpha1.MetaInfoPorts{}
	ports.GrpcLeaderElectionPort, ports.GrpcNonLeaderElectionPort = grpcregistry.GetGrpcPorts()
	newStatus := publicv1alpha1.MetaInfoStatus{
		Namespace: namespace,
		Endpoints: endpoints,
		Ports:     &ports,
	}

	metaInfo := &publicv1alpha1.MetaInfo{}
	err = r.Get(context.TODO(), request.NamespacedName, metaInfo)
	if err != nil {
		if !errors.IsNotFound(err) {
			return reconcile.Result{}, fmt.Errorf("get metainfo %s error: %v", request.Name, err)
		}

		metaInfo.Name = publicv1alpha1.MetaNameOfKruiseManager
		metaInfo.Status = newStatus
		err = r.Create(context.TODO(), metaInfo)
		if err != nil && !errors.IsAlreadyExists(err) {
			return reconcile.Result{}, fmt.Errorf("create metainfo %s error: %v", request.Name, err)
		}
		return
	}

	if reflect.DeepEqual(metaInfo.Status, newStatus) {
		return reconcile.Result{}, nil
	}

	metaInfo.Status = newStatus
	err = r.Status().Update(context.TODO(), metaInfo)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("update metainfo %s error: %v", request.Name, err)
	}
	return reconcile.Result{}, nil
}
