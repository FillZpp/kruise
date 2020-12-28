---
title: Design application safety polices
authors:
  - "@FillZpp"
reviewers:
  - "@Fei-Guo"
  - "@furykerry"
creation-date: 2020-12-28
last-updated: 2020-12-28
status: provisional
---

# Design application safety polices

## Table of Contents

A table of contents is helpful for quickly jumping to sections of a proposal and for highlighting
any additional information provided beyond the standard proposal template.
[Tools for generating](https://github.com/ekalinin/github-markdown-toc) a table of contents from markdown are available.

   * [Design application safety polices](#design-application-safety-polices)
      * [Table of Contents](#table-of-contents)
      * [Summary](#summary)
         * [Goals](#goals)
      * [Proposal](#proposal)
         * [1. Protect CRD from cascading deletion](#1-protect-crd-from-cascading-deletion)
         * [2. Protect Workload from cascading deletion](#2-protect-workload-from-cascading-deletion)
         * [3. PodUnavailableBudget (PUB)](#3-podunavailablebudget-pub)
         * [4. Generate PDB/PUB automatically for Workloads](#4-generate-pdbpub-automatically-for-workloads)
         * [5. Global flow control for Pod deletion](#5-global-flow-control-for-pod-deletion)
      * [Implementation History](#implementation-history)

## Summary

Currently, there are so many risks in a Kubernetes cluster:

1. delete a CRD mistakenly, all CR disappeared
2. delete a Workload mistakenly, all Pods belongs to it deleted
3. upgrade a Workload with incorrect updateStrategy, Pods belongs to it recreated unexpectedly
4. delete all Workloads mistakenly in batches, all applications in the cluster unavailable
5. drain lots of Nodes or patch an incorrect taint, all Pods on them evicted

We should implement safety policies to keep applications safety and prevent mistaken operations.

### Goals

- Design safety policies for cloud native applications.

## Proposal

There are several proposals for different policies.

### 1. Protect CRD from cascading deletion

Implementing a webhook for CRD deletion.

If the CRD to delete has `policy.kruise.io/disable-cascading-deletion: true` in labels or annotations and there still exists
CRs for this CRD, webhook will reject such delete operation.

### 2. Protect Workload from cascading deletion

Implementing a webhook for Workloads deletion (including Workloads in Kruise and official Kubernetes).

If the Workload to delete has `policy.kruise.io/disable-cascading-deletion: true` in labels or annotations and its replicas is not zero,
webhook will reject such delete operation.

### 3. PodUnavailableBudget (PUB)

The PodDisruptionBudget (PDB) in official Kubernetes can only limit evict operation, which is .

So we decide to design a CRD named PodUnavailableBudget, which supports to protect application from both deletion,
eviction, in-place update, readinessGate condition update and even hot-upgrade defined by users.

```yaml
apiVersion: policy.kruise.io/v1alpha1
kind: PodUnavailableBudget
spec:
  #selector:
  #  app: xxx
  targetRef:
    apiVersion: apps.kruise.io
    kind: CloneSet
    name: app-xxx
  maxUnavailable: 5
  # minAvailable: 70%
status:
  deletedPods:
    pod-uid-xxx: "116894821"
  unavailablePods:
    pod-name-xxx: "116893007"
  unavailableAllowed: 2
  currentAvailable: 17
  desiredAvailable: 15
  totalReplicas: 20
```

- `selector`/`targetRef`: selector or workload reference for the Pods to protect. The two fields are mutually exclusive.
- `maxUnavailable`/`minAvailable`: same as PDB, they can be absolute numbers or percentages. The two fields are mutually exclusive.

- `deletedPods`: contains Pods whose eviction or deletion ware processed by webhook but have not yet been observed by the PUB controller.
- `unavailablePods`: contains Pods whose updating ware processed by webhook but have not yet been observed by the PUB controller.
- `unavailableAllowed`: number of Pod unavailable that are currently allowed.
- `currentAvailable`: current number of available Pods
- `desiredAvailable`: minimum desired number of available pods
- `totalReplicas`: total number of pods counted by this PUB

### 4. Generate PDB/PUB automatically for Workloads

Since we have both PDB and PUB, Kruise can automatically generate these budgets for application workloads.

For example, users can add these annotations to a Deployment: 

```yaml
apiVersion: apps
kind: Deployment
metadata:
  name: deploy-foo
  annotations:
    policy.kruise.io/generate-pub: "true"
    policy.kruise.io/generate-pub-maxUnavailable: "20%"
    # policy.kruise.io/generate-pub-minAvailable: "80%"
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 25%
      maxSurge: 25%
  # ...
```

Then, kruise-manager will automatically create a PUB for it:

```yaml
apiVersion: policy.kruise.io/v1alpha1
kind: PodUnavailableBudget
spec:
  targetRef:
    apiVersion: apps
    kind: Deployment
    name: deploy-foo
  maxUnavailable: 20%
```

If users just set `policy.kruise.io/generate-pub=true` without other available or unavailable annotations, Kruise will take
the `maxUnavailable` defined in update strategy of this workload as the PUB/PDB strategy.

Workload types we should support:

- Official workloads in `apps` group:
  - Deployment
  - StatefulSet
  - ReplicaSet
- Kruise workloads in `apps.kruise.io` group:
  - CloneSet
  - StatefulSet
  - UnitedDeployment

### 5. Global flow control for Pod deletion

This is global flow control policy, which means it limits Pod deletion in the whole Kubernetes cluster.

```yaml
apiVersion: policy.kruise.io/v1alpha1
kind: PodDeletionFlowControl
metadata:
  # ...
spec:
  limitRules:
  - interval: 1m
    limit: 50
  - interval: 10m
    limit: 300
  - interval: 1h
    limit: 1000
  whiteListSelector:
    matchExpressions:
    - key: xxx
      operator: In
      value: foo
```

If the number of Pods deleted is bigger than the limit number of a limitRule, other Pod deletion will be rejected in this interval duration.

## Implementation History
