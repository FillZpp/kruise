/*
Copyright 2022 The Kruise Authors.

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

package selfdeletioncost

import (
	"context"
	"os"

	"github.com/openkruise/kruise/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kubeclientset "k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var (
	podNamespace = os.Getenv("POD_NAMESPACE")
	podName      = os.Getenv("POD_NAME")
)

func Add(mgr manager.Manager) error {
	if podNamespace == "" || podName == "" {
		klog.Warningf(`Skip with env POD_NAMESPACE="%s", POD_NAME="%s"`, podNamespace, podName)
		return nil
	}
	r := &setter{cli: client.GetGenericClient().KubeClient}
	r.removeDeletionCost()
	return mgr.Add(r)
}

type setter struct {
	cli kubeclientset.Interface
}

func (s *setter) NeedLeaderElection() bool {
	return true
}

func (s *setter) Start(ctx context.Context) error {
	s.addDeletionCost()
	<-ctx.Done()
	s.removeDeletionCost()
	return nil
}

func (s *setter) addDeletionCost() {
	patchBody := `{"metadata":{"annotations":{"controller.kubernetes.io/pod-deletion-cost":"100"}}}`
	_, err := s.cli.CoreV1().Pods(podNamespace).Patch(context.TODO(), podName, types.StrategicMergePatchType, []byte(patchBody), metav1.PatchOptions{})
	if err != nil {
		klog.Errorf("Failed to add deletion-cost %v to self %s/%s: %v", patchBody, podNamespace, podName, err)
	}
}

func (s *setter) removeDeletionCost() {
	patchBody := `{"metadata":{"annotations":{"controller.kubernetes.io/pod-deletion-cost":null}}}`
	_, err := s.cli.CoreV1().Pods(podNamespace).Patch(context.TODO(), podName, types.StrategicMergePatchType, []byte(patchBody), metav1.PatchOptions{})
	if err != nil {
		klog.Errorf("Failed to remove deletion-cost %v to self %s/%s: %v", patchBody, podNamespace, podName, err)
	}
}
