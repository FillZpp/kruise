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

package util

import (
	"crypto/md5"
	"encoding/hex"

	ctrlmeshproto "github.com/openkruise/kruise/apis/ctrlmesh/proto"
	"github.com/openkruise/kruise/pkg/util"
)

func CalculateHashForProtoSpec(spec *ctrlmeshproto.Spec) *ctrlmeshproto.SpecHash {
	specHash := &ctrlmeshproto.SpecHash{}
	if spec.Meta != nil {
		specHash.VopResourceVersion = spec.Meta.VopResourceVersion
	}
	if spec.Route != nil {
		specHash.RouteHash = getMD5Hash(util.DumpJSON(spec.Route))
		specHash.RouteStrictHash = getMD5Hash(util.DumpJSON(&ctrlmeshproto.Route{
			Subset:                  spec.Route.Subset,
			GlobalLimits:            spec.Route.GlobalLimits,
			GlobalExcludeNamespaces: spec.Route.GlobalExcludeNamespaces,
			SensSubsetNamespaces:    spec.Route.SensSubsetNamespaces,
		}))
	}
	if spec.Endpoints != nil {
		specHash.EndpointsHash = getMD5Hash(util.DumpJSON(spec.Endpoints))
	}
	return specHash
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
