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
// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/openkruise/kruise/apis/ctrlmesh/v1alpha1"
	scheme "github.com/openkruise/kruise/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// VirtualOperatorsGetter has a method to return a VirtualOperatorInterface.
// A group's client should implement this interface.
type VirtualOperatorsGetter interface {
	VirtualOperators(namespace string) VirtualOperatorInterface
}

// VirtualOperatorInterface has methods to work with VirtualOperator resources.
type VirtualOperatorInterface interface {
	Create(*v1alpha1.VirtualOperator) (*v1alpha1.VirtualOperator, error)
	Update(*v1alpha1.VirtualOperator) (*v1alpha1.VirtualOperator, error)
	UpdateStatus(*v1alpha1.VirtualOperator) (*v1alpha1.VirtualOperator, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.VirtualOperator, error)
	List(opts v1.ListOptions) (*v1alpha1.VirtualOperatorList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.VirtualOperator, err error)
	VirtualOperatorExpansion
}

// virtualOperators implements VirtualOperatorInterface
type virtualOperators struct {
	client rest.Interface
	ns     string
}

// newVirtualOperators returns a VirtualOperators
func newVirtualOperators(c *CtrlmeshV1alpha1Client, namespace string) *virtualOperators {
	return &virtualOperators{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the virtualOperator, and returns the corresponding virtualOperator object, and an error if there is any.
func (c *virtualOperators) Get(name string, options v1.GetOptions) (result *v1alpha1.VirtualOperator, err error) {
	result = &v1alpha1.VirtualOperator{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("virtualoperators").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of VirtualOperators that match those selectors.
func (c *virtualOperators) List(opts v1.ListOptions) (result *v1alpha1.VirtualOperatorList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.VirtualOperatorList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("virtualoperators").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested virtualOperators.
func (c *virtualOperators) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("virtualoperators").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a virtualOperator and creates it.  Returns the server's representation of the virtualOperator, and an error, if there is any.
func (c *virtualOperators) Create(virtualOperator *v1alpha1.VirtualOperator) (result *v1alpha1.VirtualOperator, err error) {
	result = &v1alpha1.VirtualOperator{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("virtualoperators").
		Body(virtualOperator).
		Do().
		Into(result)
	return
}

// Update takes the representation of a virtualOperator and updates it. Returns the server's representation of the virtualOperator, and an error, if there is any.
func (c *virtualOperators) Update(virtualOperator *v1alpha1.VirtualOperator) (result *v1alpha1.VirtualOperator, err error) {
	result = &v1alpha1.VirtualOperator{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("virtualoperators").
		Name(virtualOperator.Name).
		Body(virtualOperator).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *virtualOperators) UpdateStatus(virtualOperator *v1alpha1.VirtualOperator) (result *v1alpha1.VirtualOperator, err error) {
	result = &v1alpha1.VirtualOperator{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("virtualoperators").
		Name(virtualOperator.Name).
		SubResource("status").
		Body(virtualOperator).
		Do().
		Into(result)
	return
}

// Delete takes name of the virtualOperator and deletes it. Returns an error if one occurs.
func (c *virtualOperators) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("virtualoperators").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *virtualOperators) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("virtualoperators").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched virtualOperator.
func (c *virtualOperators) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.VirtualOperator, err error) {
	result = &v1alpha1.VirtualOperator{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("virtualoperators").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
