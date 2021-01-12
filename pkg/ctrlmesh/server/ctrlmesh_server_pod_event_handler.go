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

	podutil "k8s.io/kubernetes/pkg/api/v1/pod"

	ctrlmeshv1alpha1 "github.com/openkruise/kruise/apis/ctrlmesh/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type podEventHandler struct {
	reader client.Reader
}

func (h *podEventHandler) Create(e event.CreateEvent, q workqueue.RateLimitingInterface) {
	pod := e.Object.(*v1.Pod)
	vop, err := h.getVOPForPod(pod)
	if err != nil {
		klog.Warningf("Failed to get VOP for Pod %s/%s creation: %v", pod.Namespace, pod.Name, err)
		return
	} else if vop == nil {
		return
	}
	q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
		Namespace: vop.Namespace,
		Name:      vop.Name,
	}})
}

func (h *podEventHandler) Update(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	oldPod := e.ObjectOld.(*v1.Pod)
	newPod := e.ObjectNew.(*v1.Pod)

	if newPod.DeletionTimestamp != nil {
		return
	}

	vop, err := h.getVOPForPod(newPod)
	if err != nil {
		klog.Warningf("Failed to get VOP for Pod %s/%s update: %v", newPod.Namespace, newPod.Name, err)
		return
	} else if vop == nil {
		return
	}

	// if the subset of pod changed, enqueue it
	oldSubset := determinePodSubset(vop, oldPod)
	newSubset := determinePodSubset(vop, newPod)
	if oldSubset != newSubset || podutil.IsPodReady(oldPod) != podutil.IsPodReady(newPod) {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Namespace: vop.Namespace,
			Name:      vop.Name,
		}})
	}
}

func (h *podEventHandler) Delete(e event.DeleteEvent, q workqueue.RateLimitingInterface) {
}

func (h *podEventHandler) Generic(e event.GenericEvent, q workqueue.RateLimitingInterface) {
}

func (h *podEventHandler) getVOPForPod(pod *v1.Pod) (*ctrlmeshv1alpha1.VirtualOperator, error) {
	name := pod.Labels[ctrlmeshv1alpha1.VirtualOperatorInjectedKey]
	if name == "" {
		return nil, nil
	}
	vop := &ctrlmeshv1alpha1.VirtualOperator{}
	err := h.reader.Get(context.TODO(), types.NamespacedName{Namespace: pod.Namespace, Name: name}, vop)
	return vop, err
}
