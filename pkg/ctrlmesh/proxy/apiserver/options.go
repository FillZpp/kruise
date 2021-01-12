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

package apiserver

import (
	"fmt"
	"strings"
	"time"

	"github.com/openkruise/kruise/pkg/ctrlmesh/proxy/grpcclient"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/sets"
	utilwaitgroup "k8s.io/apimachinery/pkg/util/waitgroup"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/server"
	genericfilters "k8s.io/apiserver/pkg/server/filters"
	"k8s.io/apiserver/pkg/server/options"
	restclient "k8s.io/client-go/rest"
)

const (
	// DefaultLegacyAPIPrefix is where the legacy APIs will be located.
	DefaultLegacyAPIPrefix = "/api"

	// APIGroupPrefix is where non-legacy API group will be located.
	APIGroupPrefix = "/apis"
)

// Options contains everything necessary to create and run proxy.
type Options struct {
	Config               *restclient.Config
	SecureServingOptions *options.SecureServingOptions

	Serializer          runtime.NegotiatedSerializer
	RequestInfoResolver apirequest.RequestInfoResolver

	LegacyAPIGroupPrefixes sets.String
	LongRunningFunc        apirequest.LongRunningRequestCheck
	RequestTimeout         time.Duration
	HandlerChainWaitGroup  *utilwaitgroup.SafeWaitGroup

	MaxRequestsInFlight         int
	MaxMutatingRequestsInFlight int

	LeaderElectionName string
	GrpcClient         *grpcclient.GrpcClient
}

func NewOptions() *Options {
	o := &Options{
		Config:                 new(restclient.Config),
		SecureServingOptions:   options.NewSecureServingOptions(),
		Serializer:             newCodecs(),
		LegacyAPIGroupPrefixes: sets.NewString(DefaultLegacyAPIPrefix),
		LongRunningFunc: genericfilters.BasicLongRunningRequestCheck(
			sets.NewString("watch", "proxy"),
			sets.NewString("attach", "exec", "proxy", "log", "portforward"),
		), // BasicLongRunningRequestCheck(sets.NewString("watch"), sets.NewString()),
		RequestTimeout:              time.Duration(300) * time.Second,
		HandlerChainWaitGroup:       new(utilwaitgroup.SafeWaitGroup),
		MaxRequestsInFlight:         400,
		MaxMutatingRequestsInFlight: 800,
	}
	o.RequestInfoResolver = NewRequestInfoResolver(o)
	return o
}

func (o *Options) ApplyTo(apiserver **server.SecureServingInfo) error {
	if o == nil {
		return fmt.Errorf("SecureServingInfo is empty")
	}
	err := o.SecureServingOptions.ApplyTo(apiserver)
	if err != nil {
		return err
	}
	return err
}

func (o *Options) Validate() []error {
	errors := []error{}
	errors = append(errors, o.SecureServingOptions.Validate()...)
	return errors
}

func NewRequestInfoResolver(o *Options) *apirequest.RequestInfoFactory {
	apiPrefixes := sets.NewString(strings.Trim(APIGroupPrefix, "/")) // all possible API prefixes
	legacyAPIPrefixes := sets.String{}                               // APIPrefixes that won't have groups (legacy)
	for legacyAPIPrefix := range o.LegacyAPIGroupPrefixes {
		apiPrefixes.Insert(strings.Trim(legacyAPIPrefix, "/"))
		legacyAPIPrefixes.Insert(strings.Trim(legacyAPIPrefix, "/"))
	}

	return &apirequest.RequestInfoFactory{
		APIPrefixes:          apiPrefixes,
		GrouplessAPIPrefixes: legacyAPIPrefixes,
	}
}

func newCodecs() runtime.NegotiatedSerializer {
	scheme := runtime.NewScheme()

	admissionv1beta1.AddToScheme(scheme)
	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(scheme, schema.GroupVersion{Version: "v1"})

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)

	return serializer.NewCodecFactory(scheme)
}
