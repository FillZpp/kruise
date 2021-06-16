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

package server

import (
	"sort"
	"strconv"

	ctrlmeshutil "github.com/openkruise/kruise/pkg/ctrlmesh/util"

	ctrlmeshproto "github.com/openkruise/kruise/apis/ctrlmesh/proto"
	ctrlmeshv1alpha1 "github.com/openkruise/kruise/apis/ctrlmesh/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
)

func determinePodSubset(vop *ctrlmeshv1alpha1.VirtualOperator, pod *v1.Pod) string {
	for i := range vop.Spec.Subsets {
		subset := &vop.Spec.Subsets[i]
		matched := true
		for k, v := range subset.Labels {
			if pod.Labels[k] != v {
				matched = false
				break
			}
		}
		if matched {
			return subset.Name
		}
	}
	return ""
}

func generateProtoRoute(vop *ctrlmeshv1alpha1.VirtualOperator, namespaces []*v1.Namespace) map[string]*ctrlmeshproto.Route {
	globalLimits := vop.Spec.Route.GlobalLimits
	subRules := vop.Spec.Route.SubRules
	sort.SliceStable(subRules, func(i, j int) bool { return subRules[i].Name < subRules[j].Name })
	subsets := vop.Spec.Subsets
	sort.SliceStable(subsets, func(i, j int) bool { return subsets[i].Name < subsets[j].Name })

	routesForAllSubsets := make(map[string]*ctrlmeshproto.Route, len(vop.Spec.Subsets)+1)

	rulesMap := make(map[string][]ctrlmeshv1alpha1.MatchLimitSelector, len(vop.Spec.Route.SubRules))
	for i := range vop.Spec.Route.SubRules {
		r := &vop.Spec.Route.SubRules[i]
		rulesMap[r.Name] = r.Match
	}

	globalExcludeNamespaces := sets.NewString()
	for _, ms := range globalLimits {
		for _, ns := range namespaces {
			if globalExcludeNamespaces.Has(ns.Name) {
				continue
			}
			if match, _ := ctrlmeshutil.IsNamespaceMatchesLimitSelector(&ms, ns); match {
				globalExcludeNamespaces.Insert(ns.Name)
			}
		}
	}

	previousExcludeNamespaces := sets.NewString().Union(globalExcludeNamespaces)
	var sensSubsetNamespaces []ctrlmeshproto.InternalSensSubsetNamespaces
	// add subsets
	for i := range vop.Spec.Subsets {
		subset := &vop.Spec.Subsets[i]
		subsetNamespaces := sets.NewString()
		for _, ruleName := range subset.RouteRules {
			for _, ms := range rulesMap[ruleName] {
				for _, ns := range namespaces {
					if previousExcludeNamespaces.Has(ns.Name) {
						continue
					}
					if match, _ := ctrlmeshutil.IsNamespaceMatchesLimitSelector(&ms, ns); match {
						subsetNamespaces.Insert(ns.Name)
					}
				}
			}
		}
		sensSubsetNamespaces = append(sensSubsetNamespaces, ctrlmeshproto.InternalSensSubsetNamespaces{Name: subset.Name, Namespaces: subsetNamespaces})
		internalRoute := &ctrlmeshproto.InternalRoute{
			Subset:                  subset.Name,
			GlobalLimits:            globalLimits,
			SubRules:                subRules,
			Subsets:                 subsets,
			GlobalExcludeNamespaces: globalExcludeNamespaces,
			SensSubsetNamespaces:    sensSubsetNamespaces,
		}
		routesForAllSubsets[subset.Name] = internalRoute.Encode()
		previousExcludeNamespaces = previousExcludeNamespaces.Union(subsetNamespaces)
	}

	// add default
	routesForAllSubsets[""] = (&ctrlmeshproto.InternalRoute{
		GlobalLimits:            globalLimits,
		SubRules:                subRules,
		Subsets:                 subsets,
		GlobalExcludeNamespaces: globalExcludeNamespaces,
		SensSubsetNamespaces:    sensSubsetNamespaces,
	}).Encode()

	return routesForAllSubsets
}

func generateProtoEndpoints(vop *ctrlmeshv1alpha1.VirtualOperator, pods []*v1.Pod) []*ctrlmeshproto.Endpoint {
	var endpoints []*ctrlmeshproto.Endpoint
	for _, pod := range pods {
		if podutil.IsPodReady(pod) {
			endpoints = append(endpoints, &ctrlmeshproto.Endpoint{
				Name:   pod.Name,
				Ip:     pod.Status.PodIP,
				Subset: determinePodSubset(vop, pod),
			})
		}
	}
	if len(endpoints) > 0 {
		sort.SliceStable(endpoints, func(i, j int) bool {
			return endpoints[i].Name < endpoints[j].Name
		})
	}
	return endpoints
}

func isResourceVersionNewer(old, new string) bool {
	if len(old) == 0 {
		return true
	}

	oldCount, err := strconv.ParseUint(old, 10, 64)
	if err != nil {
		return true
	}

	newCount, err := strconv.ParseUint(new, 10, 64)
	if err != nil {
		return false
	}

	return newCount >= oldCount
}
