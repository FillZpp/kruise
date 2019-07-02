// +build !ignore_autogenerated

/*
Copyright 2019 The Kruise Authors.

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
// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.BroadcastJob":                     schema_pkg_apis_apps_v1alpha1_BroadcastJob(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.BroadcastJobList":                 schema_pkg_apis_apps_v1alpha1_BroadcastJobList(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.BroadcastJobSpec":                 schema_pkg_apis_apps_v1alpha1_BroadcastJobSpec(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.BroadcastJobStatus":               schema_pkg_apis_apps_v1alpha1_BroadcastJobStatus(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.CompletionPolicy":                 schema_pkg_apis_apps_v1alpha1_CompletionPolicy(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.InPlaceUpdateContainerStatus":     schema_pkg_apis_apps_v1alpha1_InPlaceUpdateContainerStatus(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.InPlaceUpdateState":               schema_pkg_apis_apps_v1alpha1_InPlaceUpdateState(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.JobCondition":                     schema_pkg_apis_apps_v1alpha1_JobCondition(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.RollingUpdateStatefulSetStrategy": schema_pkg_apis_apps_v1alpha1_RollingUpdateStatefulSetStrategy(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.SidecarContainer":                 schema_pkg_apis_apps_v1alpha1_SidecarContainer(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.SidecarSet":                       schema_pkg_apis_apps_v1alpha1_SidecarSet(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.SidecarSetList":                   schema_pkg_apis_apps_v1alpha1_SidecarSetList(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.SidecarSetSpec":                   schema_pkg_apis_apps_v1alpha1_SidecarSetSpec(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.SidecarSetStatus":                 schema_pkg_apis_apps_v1alpha1_SidecarSetStatus(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.StatefulSet":                      schema_pkg_apis_apps_v1alpha1_StatefulSet(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.StatefulSetList":                  schema_pkg_apis_apps_v1alpha1_StatefulSetList(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.StatefulSetSpec":                  schema_pkg_apis_apps_v1alpha1_StatefulSetSpec(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.StatefulSetStatus":                schema_pkg_apis_apps_v1alpha1_StatefulSetStatus(ref),
		"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.StatefulSetUpdateStrategy":        schema_pkg_apis_apps_v1alpha1_StatefulSetUpdateStrategy(ref),
	}
}

func schema_pkg_apis_apps_v1alpha1_BroadcastJob(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "BroadcastJob is the Schema for the broadcastjobs API",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.BroadcastJobSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.BroadcastJobStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.BroadcastJobSpec", "github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.BroadcastJobStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_apps_v1alpha1_BroadcastJobList(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "BroadcastJobList contains a list of BroadcastJob",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"),
						},
					},
					"items": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.BroadcastJob"),
									},
								},
							},
						},
					},
				},
				Required: []string{"items"},
			},
		},
		Dependencies: []string{
			"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.BroadcastJob", "k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"},
	}
}

func schema_pkg_apis_apps_v1alpha1_BroadcastJobSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "BroadcastJobSpec defines the desired state of BroadcastJob",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"parallelism": {
						SchemaProps: spec.SchemaProps{
							Description: "Specifies the maximum desired number of pods the job should run at any given time. The actual number of pods running in steady state will be less than this number when the work left to do is less than max parallelism. Not setting this value means no limit.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"template": {
						SchemaProps: spec.SchemaProps{
							Description: "Describes the pod that will be created when executing a job.",
							Ref:         ref("k8s.io/api/core/v1.PodTemplateSpec"),
						},
					},
					"completionPolicy": {
						SchemaProps: spec.SchemaProps{
							Description: "CompletionPolicy indicates the completion policy of the job. Default is Always CompletionPolicyType",
							Ref:         ref("github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.CompletionPolicy"),
						},
					},
				},
				Required: []string{"template"},
			},
		},
		Dependencies: []string{
			"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.CompletionPolicy", "k8s.io/api/core/v1.PodTemplateSpec"},
	}
}

func schema_pkg_apis_apps_v1alpha1_BroadcastJobStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "BroadcastJobStatus defines the observed state of BroadcastJob",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"conditions": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-patch-merge-key": "type",
								"x-kubernetes-patch-strategy":  "merge",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "The latest available observations of an object's current state.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.JobCondition"),
									},
								},
							},
						},
					},
					"startTime": {
						SchemaProps: spec.SchemaProps{
							Description: "Represents time when the job was acknowledged by the job controller. It is not guaranteed to be set in happens-before order across separate operations. It is represented in RFC3339 form and is in UTC.",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.Time"),
						},
					},
					"completionTime": {
						SchemaProps: spec.SchemaProps{
							Description: "Represents time when the job was completed. It is not guaranteed to be set in happens-before order across separate operations. It is represented in RFC3339 form and is in UTC.",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.Time"),
						},
					},
					"active": {
						SchemaProps: spec.SchemaProps{
							Description: "The number of actively running pods.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"succeeded": {
						SchemaProps: spec.SchemaProps{
							Description: "The number of pods which reached phase Succeeded.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"failed": {
						SchemaProps: spec.SchemaProps{
							Description: "The number of pods which reached phase Failed.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"desired": {
						SchemaProps: spec.SchemaProps{
							Description: "The desired number of pods, this is typically equal to the number of nodes satisfied to run pods.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.JobCondition", "k8s.io/apimachinery/pkg/apis/meta/v1.Time"},
	}
}

func schema_pkg_apis_apps_v1alpha1_CompletionPolicy(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "CompletionPolicy indicates the completion policy for the job",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"type": {
						SchemaProps: spec.SchemaProps{
							Description: "Type indicates the type of the CompletionPolicy Default is Always",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"activeDeadlineSeconds": {
						SchemaProps: spec.SchemaProps{
							Description: "Specifies the duration in seconds relative to the startTime that the job may be active before the system tries to terminate it; value must be positive integer. Only works for Always type",
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"backoffLimit": {
						SchemaProps: spec.SchemaProps{
							Description: "Specifies the number of retries before marking this job failed. Not setting value means no limit. Only works for Always type",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"ttlSecondsAfterFinished": {
						SchemaProps: spec.SchemaProps{
							Description: "ttlSecondsAfterFinished limits the lifetime of a Job that has finished execution (either Complete or Failed). If this field is set, ttlSecondsAfterFinished after the Job finishes, it is eligible to be automatically deleted. When the Job is being deleted, its lifecycle guarantees (e.g. finalizers) will be honored. If this field is unset, the Job won't be automatically deleted. If this field is set to zero, the Job becomes eligible to be deleted immediately after it finishes. This field is alpha-level and is only honored by servers that enable the TTLAfterFinished feature. Only works for Always type",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
				},
			},
		},
	}
}

func schema_pkg_apis_apps_v1alpha1_InPlaceUpdateContainerStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "InPlaceUpdateContainerStatus records container status in current pod.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"imageID": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
			},
		},
	}
}

func schema_pkg_apis_apps_v1alpha1_InPlaceUpdateState(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "InPlaceUpdateState records latest inplace-update state, including old statuses of containers.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"revision": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"updateTimestamp": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.Time"),
						},
					},
					"lastContainerStatuses": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.InPlaceUpdateContainerStatus"),
									},
								},
							},
						},
					},
				},
				Required: []string{"revision", "updateTimestamp", "lastContainerStatuses"},
			},
		},
		Dependencies: []string{
			"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.InPlaceUpdateContainerStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.Time"},
	}
}

func schema_pkg_apis_apps_v1alpha1_JobCondition(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "JobCondition describes current state of a job.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"type": {
						SchemaProps: spec.SchemaProps{
							Description: "Type of job condition, Complete or Failed.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Description: "Status of the condition, one of True, False, Unknown.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"lastProbeTime": {
						SchemaProps: spec.SchemaProps{
							Description: "Last time the condition was checked.",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.Time"),
						},
					},
					"lastTransitionTime": {
						SchemaProps: spec.SchemaProps{
							Description: "Last time the condition transit from one status to another.",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.Time"),
						},
					},
					"reason": {
						SchemaProps: spec.SchemaProps{
							Description: "(brief) reason for the condition's last transition.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"message": {
						SchemaProps: spec.SchemaProps{
							Description: "Human readable message indicating details about last transition.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"type", "status"},
			},
		},
		Dependencies: []string{
			"k8s.io/apimachinery/pkg/apis/meta/v1.Time"},
	}
}

func schema_pkg_apis_apps_v1alpha1_RollingUpdateStatefulSetStrategy(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "RollingUpdateStatefulSetStrategy is used to communicate parameter for RollingUpdateStatefulSetStrategyType.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"partition": {
						SchemaProps: spec.SchemaProps{
							Description: "Partition indicates the ordinal at which the StatefulSet should be partitioned. Default value is 0.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"maxUnavailable": {
						SchemaProps: spec.SchemaProps{
							Description: "The maximum number of pods that can be unavailable during the update. Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%). Absolute number is calculated from percentage by rounding down. Also, maxUnavailable can just be allowed to work with Parallel podManagementPolicy. Defaults to 1.",
							Ref:         ref("k8s.io/apimachinery/pkg/util/intstr.IntOrString"),
						},
					},
					"podUpdatePolicy": {
						SchemaProps: spec.SchemaProps{
							Description: "PodUpdatePolicy indicates how pods should be updated Default value is \"ReCreate\"",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
		Dependencies: []string{
			"k8s.io/apimachinery/pkg/util/intstr.IntOrString"},
	}
}

func schema_pkg_apis_apps_v1alpha1_SidecarContainer(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"Container": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/api/core/v1.Container"),
						},
					},
				},
				Required: []string{"Container"},
			},
		},
		Dependencies: []string{
			"k8s.io/api/core/v1.Container"},
	}
}

func schema_pkg_apis_apps_v1alpha1_SidecarSet(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "SidecarSet is the Schema for the sidecarsets API",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.SidecarSetSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.SidecarSetStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.SidecarSetSpec", "github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.SidecarSetStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_apps_v1alpha1_SidecarSetList(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "SidecarSetList contains a list of SidecarSet",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"),
						},
					},
					"items": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.SidecarSet"),
									},
								},
							},
						},
					},
				},
				Required: []string{"items"},
			},
		},
		Dependencies: []string{
			"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.SidecarSet", "k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"},
	}
}

func schema_pkg_apis_apps_v1alpha1_SidecarSetSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "SidecarSetSpec defines the desired state of SidecarSet",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"selector": {
						SchemaProps: spec.SchemaProps{
							Description: "selector is a label query over pods that should be injected",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"),
						},
					},
					"containers": {
						SchemaProps: spec.SchemaProps{
							Description: "containers contains two pieces of information: 1. normal container info that should be injected into pod 2. custom fields to control insert behavior(currently empty)",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.SidecarContainer"),
									},
								},
							},
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.SidecarContainer", "k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"},
	}
}

func schema_pkg_apis_apps_v1alpha1_SidecarSetStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "SidecarSetStatus defines the observed state of SidecarSet",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"observedGeneration": {
						SchemaProps: spec.SchemaProps{
							Description: "observedGeneration is the most recent generation observed for this SidecarSet. It corresponds to the SidecarSet's generation, which is updated on mutation by the API Server.",
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"matchedPods": {
						SchemaProps: spec.SchemaProps{
							Description: "matchedPods is the number of Pods whose labels are matched with this SidecarSet's selector",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"updatedPods": {
						SchemaProps: spec.SchemaProps{
							Description: "updatedPods is the number of matched Pods that are injected with the latest SidecarSet's containers",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"readyPods": {
						SchemaProps: spec.SchemaProps{
							Description: "readyPods is the number of matched Pods that have a ready condition",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
				},
				Required: []string{"matchedPods", "updatedPods", "readyPods"},
			},
		},
	}
}

func schema_pkg_apis_apps_v1alpha1_StatefulSet(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "StatefulSet is the Schema for the statefulsets API",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.StatefulSetSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.StatefulSetStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.StatefulSetSpec", "github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.StatefulSetStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_apps_v1alpha1_StatefulSetList(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "StatefulSetList contains a list of StatefulSet",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"),
						},
					},
					"items": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.StatefulSet"),
									},
								},
							},
						},
					},
				},
				Required: []string{"items"},
			},
		},
		Dependencies: []string{
			"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.StatefulSet", "k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"},
	}
}

func schema_pkg_apis_apps_v1alpha1_StatefulSetSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "StatefulSetSpec defines the desired state of StatefulSet",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"replicas": {
						SchemaProps: spec.SchemaProps{
							Description: "replicas is the desired number of replicas of the given Template. These are replicas in the sense that they are instantiations of the same Template, but individual replicas also have a consistent identity. If unspecified, defaults to 1.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"selector": {
						SchemaProps: spec.SchemaProps{
							Description: "selector is a label query over pods that should match the replica count. It must match the pod template's labels. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"),
						},
					},
					"template": {
						SchemaProps: spec.SchemaProps{
							Description: "template is the object that describes the pod that will be created if insufficient replicas are detected. Each pod stamped out by the StatefulSet will fulfill this Template, but have a unique identity from the rest of the StatefulSet.",
							Ref:         ref("k8s.io/api/core/v1.PodTemplateSpec"),
						},
					},
					"volumeClaimTemplates": {
						SchemaProps: spec.SchemaProps{
							Description: "volumeClaimTemplates is a list of claims that pods are allowed to reference. The StatefulSet controller is responsible for mapping network identities to claims in a way that maintains the identity of a pod. Every claim in this list must have at least one matching (by name) volumeMount in one container in the template. A claim in this list takes precedence over any volumes in the template, with the same name.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("k8s.io/api/core/v1.PersistentVolumeClaim"),
									},
								},
							},
						},
					},
					"serviceName": {
						SchemaProps: spec.SchemaProps{
							Description: "serviceName is the name of the service that governs this StatefulSet. This service must exist before the StatefulSet, and is responsible for the network identity of the set. Pods get DNS/hostnames that follow the pattern: pod-specific-string.serviceName.default.svc.cluster.local where \"pod-specific-string\" is managed by the StatefulSet controller.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"podManagementPolicy": {
						SchemaProps: spec.SchemaProps{
							Description: "podManagementPolicy controls how pods are created during initial scale up, when replacing pods on nodes, or when scaling down. The default policy is `OrderedReady`, where pods are created in increasing order (pod-0, then pod-1, etc) and the controller will wait until each pod is ready before continuing. When scaling down, the pods are removed in the opposite order. The alternative policy is `Parallel` which will create pods in parallel to match the desired scale without waiting, and on scale down will delete all pods at once.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"updateStrategy": {
						SchemaProps: spec.SchemaProps{
							Description: "updateStrategy indicates the StatefulSetUpdateStrategy that will be employed to update Pods in the StatefulSet when a revision is made to Template.",
							Ref:         ref("github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.StatefulSetUpdateStrategy"),
						},
					},
					"revisionHistoryLimit": {
						SchemaProps: spec.SchemaProps{
							Description: "revisionHistoryLimit is the maximum number of revisions that will be maintained in the StatefulSet's revision history. The revision history consists of all revisions not represented by a currently applied StatefulSetSpec version. The default value is 10.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
				},
				Required: []string{"selector", "template"},
			},
		},
		Dependencies: []string{
			"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.StatefulSetUpdateStrategy", "k8s.io/api/core/v1.PersistentVolumeClaim", "k8s.io/api/core/v1.PodTemplateSpec", "k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector"},
	}
}

func schema_pkg_apis_apps_v1alpha1_StatefulSetStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "StatefulSetStatus defines the observed state of StatefulSet",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"observedGeneration": {
						SchemaProps: spec.SchemaProps{
							Description: "observedGeneration is the most recent generation observed for this StatefulSet. It corresponds to the StatefulSet's generation, which is updated on mutation by the API Server.",
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"replicas": {
						SchemaProps: spec.SchemaProps{
							Description: "replicas is the number of Pods created by the StatefulSet controller.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"readyReplicas": {
						SchemaProps: spec.SchemaProps{
							Description: "readyReplicas is the number of Pods created by the StatefulSet controller that have a Ready Condition.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"currentReplicas": {
						SchemaProps: spec.SchemaProps{
							Description: "currentReplicas is the number of Pods created by the StatefulSet controller from the StatefulSet version indicated by currentRevision.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"updatedReplicas": {
						SchemaProps: spec.SchemaProps{
							Description: "updatedReplicas is the number of Pods created by the StatefulSet controller from the StatefulSet version indicated by updateRevision.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"currentRevision": {
						SchemaProps: spec.SchemaProps{
							Description: "currentRevision, if not empty, indicates the version of the StatefulSet used to generate Pods in the sequence [0,currentReplicas).",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"updateRevision": {
						SchemaProps: spec.SchemaProps{
							Description: "updateRevision, if not empty, indicates the version of the StatefulSet used to generate Pods in the sequence [replicas-updatedReplicas,replicas)",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"collisionCount": {
						SchemaProps: spec.SchemaProps{
							Description: "collisionCount is the count of hash collisions for the StatefulSet. The StatefulSet controller uses this field as a collision avoidance mechanism when it needs to create the name for the newest ControllerRevision.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"conditions": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-patch-merge-key": "type",
								"x-kubernetes-patch-strategy":  "merge",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Represents the latest available observations of a statefulset's current state.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("k8s.io/api/apps/v1.StatefulSetCondition"),
									},
								},
							},
						},
					},
				},
				Required: []string{"replicas", "readyReplicas", "currentReplicas", "updatedReplicas"},
			},
		},
		Dependencies: []string{
			"k8s.io/api/apps/v1.StatefulSetCondition"},
	}
}

func schema_pkg_apis_apps_v1alpha1_StatefulSetUpdateStrategy(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "StatefulSetUpdateStrategy indicates the strategy that the StatefulSet controller will use to perform updates. It includes any additional parameters necessary to perform the update for the indicated strategy.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"type": {
						SchemaProps: spec.SchemaProps{
							Description: "Type indicates the type of the StatefulSetUpdateStrategy. Default is RollingUpdate.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"rollingUpdate": {
						SchemaProps: spec.SchemaProps{
							Description: "RollingUpdate is used to communicate parameters when Type is RollingUpdateStatefulSetStrategyType.",
							Ref:         ref("github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.RollingUpdateStatefulSetStrategy"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/openkruise/kruise/pkg/apis/apps/v1alpha1.RollingUpdateStatefulSetStrategy"},
	}
}
