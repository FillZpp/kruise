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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	MetaNameOfKruiseManager = "kruise-manager"
)

// MetaInfoSpec defines the desired state of MetaInfo
type MetaInfoSpec struct {
}

// MetaInfoStatus defines the observed state of MetaInfo
type MetaInfoStatus struct {
	Namespace string            `json:"namespace,omitempty"`
	Endpoints MetaInfoEndpoints `json:"endpoints,omitempty"`
	Ports     *MetaInfoPorts    `json:"ports,omitempty"`
}

type MetaInfoEndpoints []MetaInfoEndpoint

func (e MetaInfoEndpoints) Len() int      { return len(e) }
func (e MetaInfoEndpoints) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e MetaInfoEndpoints) Less(i, j int) bool {
	return e[i].Name < e[j].Name
}

type MetaInfoEndpoint struct {
	Name   string `json:"name"`
	PodIP  string `json:"podIP"`
	Leader bool   `json:"leader"`
}

type MetaInfoPorts struct {
	GrpcLeaderElectionPort    int `json:"grpcLeaderElectionPort,omitempty"`
	GrpcNonLeaderElectionPort int `json:"grpcNonLeaderElectionPort,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:subresource:status

// MetaInfo is the Schema for the metainfoes API
type MetaInfo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MetaInfoSpec   `json:"spec,omitempty"`
	Status MetaInfoStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MetaInfoList contains a list of MetaInfo
type MetaInfoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MetaInfo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MetaInfo{}, &MetaInfoList{})
}
