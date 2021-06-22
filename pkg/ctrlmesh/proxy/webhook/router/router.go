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

package router

import (
	"fmt"
	"net/http"

	"github.com/openkruise/kruise/pkg/ctrlmesh/proxy/grpcclient"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
)

type Router interface {
	Route(req *admissionv1beta1.AdmissionRequest) (*Result, error)
}

type Result struct {
	Accept   *RouteAccept
	Redirect *RouteRedirect
	Ignore   *RouteIgnore
	Reject   *RouteReject
}

type RouteAccept struct{}
type RouteRedirect struct {
	Hosts []string
}
type RouteIgnore struct{}
type RouteReject struct {
	Code  int32
	Error error
}

func New(c *grpcclient.GrpcClient) Router {
	return &router{grpcClient: c}
}

type router struct {
	grpcClient *grpcclient.GrpcClient
}

func (r *router) Route(req *admissionv1beta1.AdmissionRequest) (*Result, error) {
	protoRoute, protoEndpoints, _ := r.grpcClient.GetProtoSpec()
	matchSubset, ok := protoRoute.DetermineNamespaceSubset(req.Namespace)
	if !ok {
		return &Result{Ignore: &RouteIgnore{}}, nil
	}

	if matchSubset == protoRoute.Subset {
		return &Result{Accept: &RouteAccept{}}, nil
	}

	var hosts []string
	for _, e := range protoEndpoints {
		if e.Subset == matchSubset {
			hosts = append(hosts, e.Ip)
		}
	}

	if len(hosts) == 0 {
		return &Result{Reject: &RouteReject{Code: http.StatusNotFound, Error: fmt.Errorf("find no endpoints for subset %s", matchSubset)}}, nil
	}

	return &Result{Redirect: &RouteRedirect{Hosts: hosts}}, nil
}
