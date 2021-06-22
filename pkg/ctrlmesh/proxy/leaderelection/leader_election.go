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

package leaderelection

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	ctrlmeshproto "github.com/openkruise/kruise/apis/ctrlmesh/proto"

	"github.com/openkruise/kruise/pkg/ctrlmesh/constants"
	"github.com/openkruise/kruise/pkg/ctrlmesh/proxy/grpcclient"
	coordinationv1 "k8s.io/api/coordination/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

var runtimeScheme = runtime.NewScheme()
var runtimeSerializer runtime.Serializer

func init() {
	utilruntime.Must(v1.AddToScheme(runtimeScheme))
	utilruntime.Must(coordinationv1.AddToScheme(runtimeScheme))
	mediaTypes := serializer.NewCodecFactory(runtimeScheme).SupportedMediaTypes()
	for _, info := range mediaTypes {
		if info.MediaType == "application/json" {
			runtimeSerializer = info.Serializer
		}
	}
}

type Wrapper interface {
	WrapLock(*request.RequestInfo, *http.Request) (bool, func(*http.Response) error, error)
}

func New(cli clientset.Interface, grpcClient *grpcclient.GrpcClient, lockName string) Wrapper {
	return &wrapper{client: cli, grpcClient: grpcClient, namespace: os.Getenv(constants.EnvPodNamespace), lockName: lockName, identities: make(map[string]*identityState)}
}

type wrapper struct {
	client     clientset.Interface
	grpcClient *grpcclient.GrpcClient
	namespace  string
	lockName   string
	identities map[string]*identityState
}

type identityState struct {
	strictHash string
	isLeader   bool
}

func (w *wrapper) WrapLock(req *request.RequestInfo, r *http.Request) (bool, func(*http.Response) error, error) {
	if !req.IsResourceRequest || req.Subresource != "" {
		return false, nil, nil
	}
	if req.Namespace != w.namespace {
		return false, nil, nil
	} else if req.Verb != "create" && req.Name != w.lockName && !strings.HasPrefix(req.Name, w.lockName+"---") {
		return false, nil, nil
	}

	var adp adapter
	gvr := schema.GroupVersionResource{Group: req.APIGroup, Version: req.APIVersion, Resource: req.Resource}
	switch gvr {
	case v1.SchemeGroupVersion.WithResource("configmaps"):
		adp = newObjectAdapter(&v1.ConfigMap{})
	case v1.SchemeGroupVersion.WithResource("endpoints"):
		adp = newObjectAdapter(&v1.Endpoints{})
	case coordinationv1.SchemeGroupVersion.WithResource("leases"):
		adp = newLeaseAdapter()
	default:
		return false, nil, nil
	}
	klog.V(5).Infof("Wrapping %s resource lock %s %s", req.Verb, req.Resource, req.Name)

	route, _, specHash := w.grpcClient.GetProtoSpec()

	switch req.Verb {
	case "create":
		if err := adp.DecodeFrom(r.Body); err != nil {
			return true, nil, err
		}

		if adp.GetName() != w.lockName {
			return false, nil, nil
		}

		lockState, err := w.checkIdentityHistory(adp, specHash)
		if err != nil {
			return true, nil, err
		}

		if route.Subset != "" {
			name := setSubsetIntoName(w.lockName, route.Subset)
			adp.SetName(name)
			r.URL.Path = strings.Replace(r.URL.Path, w.lockName, name, -1)
		}

		adp.EncodeInto(r)

		modifier := func(resp *http.Response) error {
			if resp.StatusCode == http.StatusOK {
				lockState.isLeader = true
			}
			return nil
		}
		return true, modifier, nil

	case "update":
		if err := adp.DecodeFrom(r.Body); err != nil {
			return true, nil, err
		}

		if prevSubset := getSubsetFromName(adp.GetName()); prevSubset != route.Subset {
			return true, nil, fmt.Errorf("subset changed %s -> %s", prevSubset, route.Subset)
		}

		lockState, err := w.checkIdentityHistory(adp, specHash)
		if err != nil {
			return true, nil, err
		}

		adp.EncodeInto(r)

		modifier := func(resp *http.Response) error {
			if resp.StatusCode == http.StatusOK {
				lockState.isLeader = true
			}
			return nil
		}
		return true, modifier, nil

	case "get":
		// switch the election annotations
		if route.Subset != "" {
			r.URL.Path = strings.Replace(r.URL.Path, w.lockName, setSubsetIntoName(w.lockName, route.Subset), -1)
		}
		return true, nil, nil

	default:
		klog.Infof("Ignore %s lock operation", req.Verb)
	}
	return false, nil, nil
}

func (w *wrapper) checkIdentityHistory(adp adapter, specHash *ctrlmeshproto.SpecHash) (*identityState, error) {
	holdIdentity, ok := adp.GetHoldIdentity()
	if !ok {
		return nil, fmt.Errorf("find no hold identity resource lock")
	}
	lockState, ok := w.identities[holdIdentity]
	if ok {
		if lockState.isLeader && lockState.strictHash != specHash.RouteStrictHash {
			return nil, fmt.Errorf("strict hash changed %s -> %s for leader %s", lockState.strictHash, specHash.RouteStrictHash, holdIdentity)
		}
	} else {
		lockState = &identityState{strictHash: specHash.RouteStrictHash}
		w.identities[holdIdentity] = lockState
	}
	return lockState, nil
}

func getSubsetFromName(name string) string {
	words := strings.Split(name, "---")
	if len(words) == 1 {
		return ""
	}
	return words[1]
}

func setSubsetIntoName(name, subset string) string {
	return fmt.Sprintf("%s---%s", name, subset)
}
