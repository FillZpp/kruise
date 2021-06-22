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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	VirtualOperatorInjectedKey = "ctrlmesh.kruise.io/virtual-operator-injected"
)

// VirtualOperatorSpec defines the desired state of VirtualOperator
type VirtualOperatorSpec struct {
	// Selector is a label query over pods of this Operator.
	Selector *metav1.LabelSelector `json:"selector"`
	// Configuration defines the configuration of controller and webhook in this Operator.
	Configuration VirtualOperatorConfiguration `json:"configuration"`
	// Route defines the route of this Operator including global and sub rules.
	Route VirtualOperatorRoute `json:"route,omitempty"`
	// Subsets defines the subsets for this Operator.
	Subsets []VirtualOperatorSubset `json:"subsets,omitempty"`
}

// VirtualOperatorConfiguration defines the configuration of controller or webhook of this Operator.
type VirtualOperatorConfiguration struct {
	Controller *VirtualOperatorControllerConfiguration `json:"controller,omitempty"`
	Webhook    *VirtualOperatorWebhookConfiguration    `json:"webhook,omitempty"`
}

// VirtualOperatorControllerConfiguration defines the configuration of controller in this Operator.
type VirtualOperatorControllerConfiguration struct {
	LeaderElection LeaderElectionConfig `json:"leaderElection"`
}

// LeaderElectionConfig defines configuration of leader election
type LeaderElectionConfig struct {
	LockName string `json:"lockName"`
}

// VirtualOperatorWebhookConfiguration defines the configuration of webhook in this Operator.
type VirtualOperatorWebhookConfiguration struct {
	CertDir                  string   `json:"certDir"`
	Port                     int      `json:"port"`
	ServiceName              string   `json:"serviceName,omitempty"`
	MutatingConfigurations   []string `json:"mutatingConfigurations,omitempty"`
	ValidatingConfigurations []string `json:"validatingConfigurations,omitempty"`
}

// VirtualOperatorRoute defines the route of this Operator including global and sub rules.
type VirtualOperatorRoute struct {
	GlobalLimits []MatchLimitSelector          `json:"globalLimits,omitempty"`
	SubRules     []VirtualOperatorRouteSubRule `json:"subRules,omitempty"`
}

type VirtualOperatorRouteSubRule struct {
	Name  string               `json:"name"`
	Match []MatchLimitSelector `json:"match"`
}

type MatchLimitSelector struct {
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty"`
	NamespaceRegex    *string               `json:"namespaceRegex,omitempty"`
	// TODO(FillZpp): should we support objectSelector?
	//ObjectSelector    *metav1.LabelSelector `json:"objectSelector,omitempty"`
}

type VirtualOperatorSubset struct {
	Name       string            `json:"name"`
	Labels     map[string]string `json:"labels"`
	RouteRules []string          `json:"routeRules"`
}

// VirtualOperatorStatus defines the observed state of VirtualOperator
type VirtualOperatorStatus struct {
}

// +genclient
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=vop

// VirtualOperator is the Schema for the virtualoperators API
type VirtualOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualOperatorSpec   `json:"spec,omitempty"`
	Status VirtualOperatorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VirtualOperatorList contains a list of VirtualOperator
type VirtualOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualOperator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualOperator{}, &VirtualOperatorList{})
}
