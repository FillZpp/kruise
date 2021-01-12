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

package mutating

import (
	"context"
	"encoding/json"
	"net/http"

	ctrlmeshv1alpha1 "github.com/openkruise/kruise/apis/ctrlmesh/v1alpha1"
	"github.com/openkruise/kruise/pkg/util"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// VirtualOperatorCreateUpdateHandler handles VirtualOperator
type VirtualOperatorCreateUpdateHandler struct {
	// Decoder decodes objects
	Decoder *admission.Decoder
}

var _ admission.Handler = &VirtualOperatorCreateUpdateHandler{}

// Handle handles admission requests.
func (h *VirtualOperatorCreateUpdateHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	obj := &ctrlmeshv1alpha1.VirtualOperator{}
	err := h.Decoder.Decode(req, obj)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// TODO(FillZpp): set defaults

	marshalled, err := json.Marshal(obj)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	resp := admission.PatchResponseFromRaw(req.AdmissionRequest.Object.Raw, marshalled)
	if len(resp.Patches) > 0 {
		klog.V(5).Infof("Admit VirtualOperator %s patches: %v", obj.Name, util.DumpJSON(resp.Patches))
	}
	return resp
}

var _ admission.DecoderInjector = &VirtualOperatorCreateUpdateHandler{}

// InjectDecoder injects the decoder into the VirtualOperatorCreateUpdateHandler
func (h *VirtualOperatorCreateUpdateHandler) InjectDecoder(d *admission.Decoder) error {
	h.Decoder = d
	return nil
}
