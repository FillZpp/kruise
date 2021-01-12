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
	"reflect"
	"time"

	ctrlmeshutil "github.com/openkruise/kruise/pkg/ctrlmesh/util"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	ctrlmeshv1alpha1 "github.com/openkruise/kruise/apis/ctrlmesh/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type namespaceEventHandler struct {
	reader client.Reader
}

func (h *namespaceEventHandler) Create(e event.CreateEvent, q workqueue.RateLimitingInterface) {
	vops := h.getSensitiveVOPs()
	for _, vop := range vops {
		h.enqueue(vop, q)
	}
}

func (h *namespaceEventHandler) Update(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	newNS := e.ObjectNew.(*v1.Namespace)
	oldNS := e.ObjectOld.(*v1.Namespace)
	if reflect.DeepEqual(newNS.Labels, oldNS.Labels) {
		return
	}

	vops := h.getSensitiveVOPs()
	for _, vop := range vops {
		var diff bool
		for _, ms := range vop.Spec.Route.GlobalLimits {
			oldMatch, _ := ctrlmeshutil.IsNamespaceMatchesLimitSelector(&ms, oldNS)
			newMatch, _ := ctrlmeshutil.IsNamespaceMatchesLimitSelector(&ms, newNS)
			if oldMatch != newMatch {
				diff = true
				break
			}
		}
		if diff {
			h.enqueue(vop, q)
			continue
		}

		for _, r := range vop.Spec.Route.SubRules {
			for _, ms := range r.Match {
				oldMatch, _ := ctrlmeshutil.IsNamespaceMatchesLimitSelector(&ms, oldNS)
				newMatch, _ := ctrlmeshutil.IsNamespaceMatchesLimitSelector(&ms, newNS)
				if oldMatch != newMatch {
					diff = true
					break
				}
			}
			if diff {
				break
			}
		}
		if diff {
			h.enqueue(vop, q)
		}
	}
}

func (h *namespaceEventHandler) Delete(e event.DeleteEvent, q workqueue.RateLimitingInterface) {
	vops := h.getSensitiveVOPs()
	for _, vop := range vops {
		h.enqueue(vop, q)
	}
}

func (h *namespaceEventHandler) Generic(e event.GenericEvent, q workqueue.RateLimitingInterface) {
}

func (h *namespaceEventHandler) getSensitiveVOPs() []*ctrlmeshv1alpha1.VirtualOperator {
	vopList := ctrlmeshv1alpha1.VirtualOperatorList{}
	if err := h.reader.List(context.TODO(), &vopList); err != nil {
		klog.Errorf("Failed to list all VOPs: %v", err)
		return nil
	}
	var vops []*ctrlmeshv1alpha1.VirtualOperator
	for i := range vopList.Items {
		vop := &vopList.Items[i]
		if len(vop.Spec.Route.GlobalLimits) == 0 && len(vop.Spec.Route.SubRules) == 0 {
			continue
		}
		vops = append(vops, vop)
	}
	return vops
}

func (h *namespaceEventHandler) enqueue(vop *ctrlmeshv1alpha1.VirtualOperator, q workqueue.RateLimitingInterface) {
	q.AddAfter(reconcile.Request{NamespacedName: types.NamespacedName{Namespace: vop.Namespace, Name: vop.Name}}, time.Millisecond*200)
}
