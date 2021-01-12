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
	"flag"
	"fmt"
	"strconv"
	"strings"

	ctrlmeshv1alpha1 "github.com/openkruise/kruise/apis/ctrlmesh/v1alpha1"
	"github.com/openkruise/kruise/pkg/ctrlmesh/constants"
	"github.com/openkruise/kruise/pkg/util"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/klog"
	utilpointer "k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	initImage  = flag.String("ctrlmesh-init-image", "", "The image for ControllerMesh init container.")
	proxyImage = flag.String("ctrlmesh-proxy-image", "", "The image for ControllerMesh proxy container.")

	proxyResourceCPU    = flag.String("ctrlmesh-proxy-cpu", "100m", "The CPU limit for ControllerMesh proxy container.")
	proxyResourceMemory = flag.String("ctrlmesh-proxy-memory", "200Mi", "The Memory limit for ControllerMesh proxy container.")
	proxyLogLevel       = flag.Uint("ctrlmesh-proxy-logv", 4, "The log level of ControllerMesh proxy container.")
	proxyExtraEnvs      = flag.String("ctrlmesh-extra-envs", "", "Extra environments for ControllerMesh proxy container.")
)

// +kubebuilder:rbac:groups=ctrlmesh.kruise.io,resources=virtualoperators,verbs=get;list;watch

func (h *PodCreateHandler) virtualoperatorMutatingPod(ctx context.Context, req admission.Request, pod *v1.Pod) (retErr error) {
	if req.Operation != admissionv1beta1.Create {
		return
	}

	virtualOperatorList := &ctrlmeshv1alpha1.VirtualOperatorList{}
	if err := h.Client.List(ctx, virtualOperatorList, client.InNamespace(pod.Namespace)); err != nil {
		return err
	}

	var matchedVOP *ctrlmeshv1alpha1.VirtualOperator
	for i := range virtualOperatorList.Items {
		vop := &virtualOperatorList.Items[i]
		selector, err := metav1.LabelSelectorAsSelector(vop.Spec.Selector)
		if err != nil {
			klog.Warningf("Failed to convert selector for VirtualOperator %s/%s: %v", vop.Namespace, vop.Name, err)
			continue
		}
		if selector.Matches(labels.Set(pod.Labels)) {
			if matchedVOP != nil {
				klog.Warningf("Find multiple VirtualOperator %s %s matched Pod %s/%s", matchedVOP.Name, vop.Name, pod.Namespace, pod.Name)
				return fmt.Errorf("multiple VirtualOperator %s %s matched", matchedVOP.Name, vop.Name)
			}
			matchedVOP = vop
		}
	}
	if matchedVOP == nil {
		return nil
	}

	var initContainer *v1.Container
	var proxyContainer *v1.Container
	defer func() {
		if retErr == nil {
			klog.Infof("Successfully inject VirtualOperator %s for Pod %s/%s creation, init: %s, sidecar: %s",
				matchedVOP.Name, pod.Namespace, pod.Name, util.DumpJSON(initContainer), util.DumpJSON(proxyContainer))
		} else {
			klog.Warningf("Failed to inject VirtualOperator %s for Pod %s/%s creation, error: %v",
				matchedVOP.Name, pod.Namespace, pod.Name, retErr)
		}
	}()

	if pod.Spec.HostNetwork {
		return fmt.Errorf("can not use ControllerMesh for Pod with host network")
	}
	if *initImage == "" || *proxyImage == "" {
		return fmt.Errorf("the images for ControllerMesh init or proxy container have not set in args")
	}

	initContainer = &v1.Container{
		Name:            constants.InitContainerName,
		Image:           *initImage,
		ImagePullPolicy: v1.PullAlways,
		SecurityContext: &v1.SecurityContext{
			Privileged:   utilpointer.BoolPtr(true),
			Capabilities: &v1.Capabilities{Add: []v1.Capability{"NET_ADMIN"}},
		},
	}
	proxyContainer = &v1.Container{
		Name:            constants.ProxyContainerName,
		Image:           *proxyImage,
		ImagePullPolicy: v1.PullAlways,
		Args: []string{
			"--v=" + strconv.Itoa(int(*proxyLogLevel)),
		},
		Env: []v1.EnvVar{
			{Name: constants.EnvPodName, ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{FieldPath: "metadata.name"}}},
			{Name: constants.EnvPodNamespace, ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{FieldPath: "metadata.namespace"}}},
			{Name: constants.EnvPodIP, ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{FieldPath: "status.podIP"}}},
		},
		Lifecycle: &v1.Lifecycle{
			PostStart: &v1.Handler{
				Exec: &v1.ExecAction{Command: []string{"/bin/sh", "-c", "/poststart.sh"}},
			},
		},
		ReadinessProbe: &v1.Probe{
			Handler:       v1.Handler{HTTPGet: &v1.HTTPGetAction{Path: "/healthz", Port: intstr.FromInt(constants.ProxyMetricsHealthPort)}},
			PeriodSeconds: 3,
		},
		Resources: v1.ResourceRequirements{
			Limits: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse(*proxyResourceCPU),
				v1.ResourceMemory: resource.MustParse(*proxyResourceMemory),
			},
			Requests: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("0"),
				v1.ResourceMemory: resource.MustParse("0"),
			},
		},
		SecurityContext: &v1.SecurityContext{
			Privileged:             utilpointer.BoolPtr(true), // This can be false, but true help us debug more easier.
			ReadOnlyRootFilesystem: utilpointer.BoolPtr(true),
			RunAsUser:              utilpointer.Int64Ptr(int64(constants.ProxyUserID)),
		},
	}

	if envs := getExtraEnvs(); len(envs) > 0 {
		proxyContainer.Env = append(proxyContainer.Env, envs...)
	}

	apiserverHostPortEnvs, err := getKubernetesServiceHostPort(pod)
	if err != nil {
		return err
	}
	if len(apiserverHostPortEnvs) > 0 {
		initContainer.Env = append(initContainer.Env, apiserverHostPortEnvs...)
		proxyContainer.Env = append(proxyContainer.Env, apiserverHostPortEnvs...)
	}

	if matchedVOP.Spec.Configuration.Controller != nil {
		proxyContainer.Args = append(
			proxyContainer.Args,
			fmt.Sprintf("--%s=%v", constants.ProxyLeaderElectionNameFlag, matchedVOP.Spec.Configuration.Controller.LeaderElection.LockName),
		)
	}

	if matchedVOP.Spec.Configuration.Webhook != nil {
		initContainer.Env = append(
			initContainer.Env,
			v1.EnvVar{Name: constants.EnvInboundWebhookPort, Value: strconv.Itoa(matchedVOP.Spec.Configuration.Webhook.Port)},
		)
		proxyContainer.Args = append(
			proxyContainer.Args,
			fmt.Sprintf("--%s=%v", constants.ProxyWebhookCertDirFlag, matchedVOP.Spec.Configuration.Webhook.CertDir),
			fmt.Sprintf("--%s=%v", constants.ProxyWebhookServePortFlag, matchedVOP.Spec.Configuration.Webhook.Port),
		)

		certVolumeMounts := getCertVolumeMounts(pod, matchedVOP.Spec.Configuration.Webhook.CertDir)
		if len(certVolumeMounts) > 1 {
			return fmt.Errorf("find multiple volume mounts that mount at %s: %v", matchedVOP.Spec.Configuration.Webhook.CertDir, certVolumeMounts)
		} else if len(certVolumeMounts) == 0 {
			return fmt.Errorf("find no volume mounts that mount at %s", matchedVOP.Spec.Configuration.Webhook.CertDir)
		}
		proxyContainer.VolumeMounts = append(proxyContainer.VolumeMounts, certVolumeMounts[0])
	}

	pod.Spec.InitContainers = append([]v1.Container{*initContainer}, pod.Spec.InitContainers...)
	pod.Spec.Containers = append([]v1.Container{*proxyContainer}, pod.Spec.Containers...)
	if pod.Labels == nil {
		pod.Labels = map[string]string{}
	}
	pod.Labels[ctrlmeshv1alpha1.VirtualOperatorInjectedKey] = matchedVOP.Name
	return nil
}

func getKubernetesServiceHostPort(pod *v1.Pod) (vars []v1.EnvVar, err error) {
	var hostEnv *v1.EnvVar
	var portEnv *v1.EnvVar
	for i := range pod.Spec.Containers {
		if envVar := util.GetContainerEnvVar(&pod.Spec.Containers[i], "KUBERNETES_SERVICE_HOST"); envVar != nil {
			if hostEnv != nil && hostEnv.Value != envVar.Value {
				return nil, fmt.Errorf("found multiple KUBERNETES_SERVICE_HOST values: %v, %v", hostEnv.Value, envVar.Value)
			}
			hostEnv = envVar
		}
		if envVar := util.GetContainerEnvVar(&pod.Spec.Containers[i], "KUBERNETES_SERVICE_PORT"); envVar != nil {
			if portEnv != nil && portEnv.Value != envVar.Value {
				return nil, fmt.Errorf("found multiple KUBERNETES_SERVICE_PORT values: %v, %v", portEnv.Value, envVar.Value)
			}
			portEnv = envVar
		}
	}
	if hostEnv != nil {
		vars = append(vars, *hostEnv)
	}
	if portEnv != nil {
		vars = append(vars, *portEnv)
	}
	return vars, nil
}

func getCertVolumeMounts(pod *v1.Pod, cerDir string) (mounts []v1.VolumeMount) {
	for i := range pod.Spec.Containers {
		c := &pod.Spec.Containers[i]

		for _, m := range c.VolumeMounts {
			if m.MountPath == cerDir {
				// do not modify the ref
				m.ReadOnly = true
				mounts = append(mounts, m)
			}
		}
	}
	return
}

func getExtraEnvs() (envs []v1.EnvVar) {
	if len(*proxyExtraEnvs) == 0 {
		return
	}
	kvs := strings.Split(*proxyExtraEnvs, ";")
	for _, str := range kvs {
		kv := strings.Split(str, "=")
		if len(kv) != 2 {
			continue
		}
		envs = append(envs, v1.EnvVar{Name: kv[0], Value: kv[1]})
	}
	return
}
