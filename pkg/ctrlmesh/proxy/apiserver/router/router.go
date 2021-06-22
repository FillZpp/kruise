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
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"k8s.io/apimachinery/pkg/runtime"

	ctrlmeshproto "github.com/openkruise/kruise/apis/ctrlmesh/proto"

	"github.com/openkruise/kruise/pkg/ctrlmesh/proxy/grpcclient"
	"github.com/openkruise/kruise/pkg/util"
	utildiscovery "github.com/openkruise/kruise/pkg/util/discovery"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/klog"
)

type Router interface {
	Route(*http.Request, *request.RequestInfo) (*Result, error)
}

type router struct {
	grpcClient *grpcclient.GrpcClient
}

type Result struct {
	Accept *RouteAccept
	Reject *RouteReject
}

type RouteAccept struct {
	ModifyResponse func(response *http.Response) error
	ModifyBody     func(*http.Response) io.Reader
}

type RouteReject struct {
	Code  int
	Error string
}

func New(c *grpcclient.GrpcClient) Router {
	return &router{grpcClient: c}
}

func (r *router) Route(httpReq *http.Request, reqInfo *request.RequestInfo) (*Result, error) {
	if !reqInfo.IsResourceRequest {
		return &Result{Accept: &RouteAccept{}}, nil
	}

	gvr := schema.GroupVersionResource{Group: reqInfo.APIGroup, Version: reqInfo.APIVersion, Resource: reqInfo.Resource}
	apiResource, err := utildiscovery.DiscoverGVR(gvr)
	if err != nil {
		// let requests with non-existing resources go
		if errors.IsNotFound(err) {
			return &Result{Accept: &RouteAccept{}}, nil
		}
		return nil, fmt.Errorf("failed to get gvr %v from discovery: %v", gvr, err)
	}

	tc := &transformerConfig{
		grpcClient:  r.grpcClient,
		httpReq:     httpReq,
		reqInfo:     reqInfo,
		apiResource: apiResource,
	}
	switch reqInfo.Verb {
	case "list":
		transformer := &listTransformer{transformerConfig: *tc}
		return &Result{Accept: &RouteAccept{ModifyResponse: transformer.transform}}, nil

	case "watch":
		transformer := &watchTransformer{transformerConfig: *tc}
		return &Result{Accept: &RouteAccept{ModifyBody: transformer.transform}}, nil

	default:
	}

	protoRoute, _, _ := r.grpcClient.GetProtoSpec()
	if !protoRoute.IsNamespaceMatch(reqInfo.Namespace) {
		return &Result{Reject: &RouteReject{Code: http.StatusNotFound, Error: "not match subset rules"}}, nil
	}

	return &Result{Accept: &RouteAccept{}}, nil
}

type transformerConfig struct {
	grpcClient  *grpcclient.GrpcClient
	httpReq     *http.Request
	reqInfo     *request.RequestInfo
	apiResource *metav1.APIResource
	serializers *responseSerializer
}

type listTransformer struct {
	transformerConfig
}

func (t *listTransformer) transform(resp *http.Response) error {
	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusPartialContent {
		return nil
	}
	protoRoute, _, _ := t.grpcClient.GetProtoSpec()
	// TODO: enable these lines
	//if protoRoute.Subset == "" && protoRoute.ExcludeNamespaces.Len() == 0 {
	//	return nil
	//}

	if resp.ContentLength > hugeBufferBytes {
		if err := hugeListSemaphore.Acquire(t.httpReq.Context(), 1); err != nil {
			return fmt.Errorf("acquire huge list semaphore error: %v", err)
		}
		defer hugeListSemaphore.Release(1)
	}

	t.serializers = newResponseSerializer(resp, t.apiResource, false)
	readerCloser, length, err := t.read(protoRoute)
	if err != nil {
		return err
	}
	resp.Body = readerCloser
	resp.Header.Set("Content-Length", strconv.Itoa(length))
	resp.ContentLength = int64(length)
	return nil
}

func (t *listTransformer) read(protoRoute *ctrlmeshproto.InternalRoute) (io.ReadCloser, int, error) {
	obj, err := t.serializers.DecodeList()
	if err != nil {
		t.serializers.Release()
		klog.Errorf("ListTransformer failed to decode during read list: %v", err)
		return nil, 0, err
	}

	if obj.IsList() {
		t.filterItems(protoRoute, obj)
	} else {
		klog.Warningf("ListTransformer read list is not list: %v", util.DumpJSON(t.reqInfo))
	}

	readerCloser, length, err := t.serializers.EncodeList(obj)
	if err != nil {
		t.serializers.Release()
		klog.Errorf("Transformer failed to encode during read list: %v", err)
		return nil, 0, err
	}
	return readerCloser, length, nil
}

func (t *listTransformer) filterItems(protoRoute *ctrlmeshproto.InternalRoute, listObj *unstructured.Unstructured) {
	newItems := make([]unstructured.Unstructured, 0)
	_ = listObj.EachListItem(func(obj runtime.Object) error {
		unstructuredObj := obj.(*unstructured.Unstructured)
		if protoRoute.IsNamespaceMatch(unstructuredObj.GetNamespace()) {
			newItems = append(newItems, *unstructuredObj)
		} else {
			klog.V(5).Infof("ListTransformer filter item %s/%s", unstructuredObj.GetNamespace(), unstructuredObj.GetName())
		}
		return nil
	})

	listObj.Object["items"] = newItems
	return
}

type watchTransformer struct {
	transformerConfig
	triggerChan chan struct{}
	remaining   []byte
}

func (t *watchTransformer) transform(resp *http.Response) io.Reader {
	if resp.StatusCode != http.StatusOK {
		return resp.Body
	}
	t.serializers = newResponseSerializer(resp, t.apiResource, true)
	t.triggerChan = make(chan struct{}, 1)
	return t
}

func (t *watchTransformer) Read(p []byte) (int, error) {
	// Return whatever remaining data exists from an in progress frame
	if n := len(t.remaining); n > 0 {
		if n <= len(p) {
			p = append(p[0:0], t.remaining...)
			t.remaining = nil
			return n, nil
		}

		n = len(p)
		p = append(p[0:0], t.remaining[:n]...)
		t.remaining = t.remaining[n:]
		return n, nil
	}

	body, err := t.read()
	if err != nil {
		t.serializers.Release()
		return 0, err
	}

	n := len(p)
	// If capacity of data is less than length of the message, decoder will allocate a new slice
	// and set m to it, which means we need to copy the partial result back into data and preserve
	// the remaining result for subsequent reads.
	if len(body) > n {
		p = append(p[0:0], body[:n]...)
		t.remaining = body[n:]
		return n, nil
	}
	p = append(p[0:0], body...)
	return len(body), nil
}

func (t *watchTransformer) read() ([]byte, error) {
	var err error
	var e *metav1.WatchEvent
	go func() {
		e, err = t.serializers.DecodeWatch()
		t.triggerChan <- struct{}{}
	}()

	ctx := t.httpReq.Context()
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	case <-t.triggerChan:
	}

	if err != nil {
		klog.Errorf("WatchTransformer failed to decode during read watch: %v", err)
		return nil, err
	}

	t.filterEvent(e)

	body, err := t.serializers.EncodeWatch(e)
	if err != nil {
		klog.Errorf("WatchTransformer failed to encode during read watch: %v", err)
		return nil, err
	}
	return body, nil
}

func (t *watchTransformer) filterEvent(e *metav1.WatchEvent) {
	meta, err := meta.Accessor(e.Object.Object)
	if err != nil {
		return
	}

	protoRoute, _, _ := t.grpcClient.GetProtoSpec()
	if !protoRoute.IsNamespaceMatch(meta.GetNamespace()) {
		e.Type = string(watch.Bookmark)
		klog.V(5).Infof("WatchTransformer filter object %s/%s", meta.GetNamespace(), meta.GetName())
	}
}
