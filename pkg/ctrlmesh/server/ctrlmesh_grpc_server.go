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
	"fmt"
	"io"
	"reflect"
	"sync"

	ctrlmeshproto "github.com/openkruise/kruise/apis/ctrlmesh/proto"
	ctrlmeshv1alpha1 "github.com/openkruise/kruise/apis/ctrlmesh/v1alpha1"
	"github.com/openkruise/kruise/pkg/grpcregistry"
	"github.com/openkruise/kruise/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	kubecontroller "k8s.io/kubernetes/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

var (
	grpcServer = &CtrlMeshGrpcServer{}

	grpcRecvTriggerChannel = make(chan event.GenericEvent, 1024)

	// cachedGrpcSrvConnection type is map[types.NamespacedName]*grpcSrvConnection
	cachedGrpcSrvConnection = &sync.Map{}

	// expectationsSrvHash type is map[types.NamespacedName]*ctrlmeshproto.SpecHash
	expectationsSrvHash = &sync.Map{}
)

type grpcSrvConnection struct {
	srv      ctrlmeshproto.ControllerMesh_RegisterServer
	stopChan chan struct{}

	status grpcSrvStatus
	mu     sync.Mutex
}

type grpcSrvStatus struct {
	specHash *ctrlmeshproto.SpecHash
}

func init() {
	_ = grpcregistry.Register("ctrlmesh-server", true, func(opts grpcregistry.RegisterOptions) {
		grpcServer.reader = opts.Mgr.GetCache()
		grpcServer.stopChan = opts.StopChan
		ctrlmeshproto.RegisterControllerMeshServer(opts.GrpcServer, grpcServer)
	})
}

type CtrlMeshGrpcServer struct {
	reader   client.Reader
	stopChan <-chan struct{}
}

var _ ctrlmeshproto.ControllerMeshServer = &CtrlMeshGrpcServer{}

func (s *CtrlMeshGrpcServer) Register(srv ctrlmeshproto.ControllerMesh_RegisterServer) error {
	// receive the first register message
	pStatus, err := srv.Recv()
	if err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}
	if pStatus.SelfInfo == nil || pStatus.SelfInfo.Namespace == "" || pStatus.SelfInfo.Name == "" {
		return status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid selfInfo: %+v", pStatus.SelfInfo))
	}

	// get pod
	podNamespacedName := types.NamespacedName{Namespace: pStatus.SelfInfo.Namespace, Name: pStatus.SelfInfo.Name}
	pod := &v1.Pod{}
	if err := s.reader.Get(context.TODO(), podNamespacedName, pod); err != nil {
		if errors.IsNotFound(err) {
			return status.Errorf(codes.NotFound, fmt.Sprintf("not found pod %s", podNamespacedName))
		}
		return status.Errorf(codes.Internal, fmt.Sprintf("get pod %s error: %v", podNamespacedName, err))
	} else if !kubecontroller.IsPodActive(pod) {
		return status.Errorf(codes.Canceled, fmt.Sprintf("find pod %s inactive", podNamespacedName))
	}
	vopName := pod.Labels[ctrlmeshv1alpha1.VirtualOperatorInjectedKey]
	if vopName == "" {
		return status.Errorf(codes.InvalidArgument, fmt.Sprintf("empty %s label in pod %s", ctrlmeshv1alpha1.VirtualOperatorInjectedKey, podNamespacedName))
	}

	klog.V(3).Infof("Start proxy connection from Pod %s in VOP %s", podNamespacedName, vopName)

	stopChan := make(chan struct{})
	conn := &grpcSrvConnection{srv: srv, stopChan: stopChan, status: grpcSrvStatus{specHash: pStatus.SpecHash}}
	cachedGrpcSrvConnection.Store(podNamespacedName, conn)
	expectationsSrvHash.Delete(podNamespacedName)

	grpcRecvTriggerChannel <- event.GenericEvent{Meta: &metav1.ObjectMeta{Namespace: podNamespacedName.Namespace, Name: vopName}}
	go func() {
		for {
			pStatus, err = srv.Recv()
			if err != nil {
				if err == io.EOF {
					return
				}
				select {
				case <-srv.Context().Done():
				default:
					klog.Errorf("Receive error from Pod %s in VOP %s: %v", podNamespacedName, vopName, err)
				}
				return
			}
			klog.Infof("Get proto status from Pod %s in VOP %s: %v", podNamespacedName, vopName, util.DumpJSON(pStatus))

			if pStatus.SpecHash != nil {
				var trigger bool
				conn.mu.Lock()
				if !reflect.DeepEqual(conn.status.specHash, pStatus.SpecHash) {
					conn.status = grpcSrvStatus{specHash: pStatus.SpecHash}
					trigger = true
				}
				conn.mu.Unlock()
				if v, ok := expectationsSrvHash.Load(podNamespacedName); ok {
					expectSpecHash := v.(*ctrlmeshproto.SpecHash)
					if isResourceVersionNewer(expectSpecHash.VopResourceVersion, pStatus.SpecHash.VopResourceVersion) {
						expectationsSrvHash.Delete(podNamespacedName)
					}
				}
				if trigger {
					grpcRecvTriggerChannel <- event.GenericEvent{Meta: &metav1.ObjectMeta{Namespace: podNamespacedName.Namespace, Name: vopName}}
				}
			}
		}
	}()

	select {
	case <-s.stopChan:
	case <-stopChan:
	case <-srv.Context().Done():
	}
	cachedGrpcSrvConnection.Delete(podNamespacedName)
	grpcRecvTriggerChannel <- event.GenericEvent{Meta: &metav1.ObjectMeta{Namespace: podNamespacedName.Namespace, Name: vopName}}
	klog.V(3).Infof("Finish proxy connection from Pod %s in VOP %s", podNamespacedName, vopName)
	return nil
}
