/*
Copyright 2020 The Kruise Authors.

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

package specifieddelete

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
	"github.com/openkruise/kruise/pkg/util"
)

func IsSpecifiedDelete(obj metav1.Object) bool {
	_, ok := obj.GetLabels()[appsv1alpha1.SpecifiedDeleteKey]
	return ok
}

func ShouldKeepPVC(obj metav1.Object) bool {
	return obj.GetLabels()[appsv1alpha1.KeepPVCForDeletionKey] == "true"
}

func PatchPodSpecifiedDelete(c client.Client, pod *v1.Pod, keepPVC bool) (bool, error) {
	if _, ok := pod.Labels[appsv1alpha1.SpecifiedDeleteKey]; ok {
		return false, nil
	}

	body := patchBody{Metadata: patchMeta{Labels: map[string]string{
		appsv1alpha1.SpecifiedDeleteKey: "true",
	}}}
	if keepPVC {
		body.Metadata.Labels[appsv1alpha1.KeepPVCForDeletionKey] = "true"
	}
	return true, c.Patch(context.TODO(), pod, client.RawPatch(types.StrategicMergePatchType, []byte(util.DumpJSON(body))))
}

type patchBody struct {
	Metadata patchMeta `json:"metadata"`
}

type patchMeta struct {
	Labels map[string]string `json:"labels"`
}
