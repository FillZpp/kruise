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

package grpcclient

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"

	ctrlmeshproto "github.com/openkruise/kruise/apis/ctrlmesh/proto"
	publicv1alpha1 "github.com/openkruise/kruise/apis/public/v1alpha1"
	"github.com/openkruise/kruise/pkg/client"
	publicv1alpha1informer "github.com/openkruise/kruise/pkg/client/informers/externalversions/public/v1alpha1"
	publicv1alpha1lister "github.com/openkruise/kruise/pkg/client/listers/public/v1alpha1"
	"github.com/openkruise/kruise/pkg/ctrlmesh/constants"
	ctrlmeshutil "github.com/openkruise/kruise/pkg/ctrlmesh/util"
	"github.com/openkruise/kruise/pkg/util"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

var (
	selfInfo = &ctrlmeshproto.SelfInfo{Namespace: os.Getenv(constants.EnvPodNamespace), Name: os.Getenv(constants.EnvPodName)}
	onceInit sync.Once
)

type GrpcClient struct {
	informer cache.SharedIndexInformer
	lister   publicv1alpha1lister.MetaInfoLister

	route     *ctrlmeshproto.InternalRoute
	endpoints []*ctrlmeshproto.Endpoint
	specHash  *ctrlmeshproto.SpecHash
	mu        sync.RWMutex
}

func New() *GrpcClient {
	return &GrpcClient{}
}

func (c *GrpcClient) Start(stopChan <-chan struct{}) error {
	clientset := client.GetGenericClient().KruiseClient
	c.informer = publicv1alpha1informer.NewFilteredMetaInfoInformer(clientset, 0, cache.Indexers{}, func(opts *metav1.ListOptions) {
		opts.FieldSelector = "metadata.name=" + publicv1alpha1.MetaNameOfKruiseManager
	})
	c.lister = publicv1alpha1lister.NewMetaInfoLister(c.informer.GetIndexer())

	go func() {
		c.informer.Run(stopChan)
	}()
	if ok := cache.WaitForCacheSync(stopChan, c.informer.HasSynced); !ok {
		return fmt.Errorf("wait metainfo informer synced error")
	}

	initChan := make(chan struct{})
	go c.connect(stopChan, initChan)
	<-initChan
	return nil
}

func (c *GrpcClient) connect(stopChan <-chan struct{}, initChan chan struct{}) {
	parentCtx, parentCancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-stopChan:
			parentCancel()
		}
	}()
	for i := 0; ; i++ {
		if i > 0 {
			time.Sleep(time.Second * 3)
		}
		klog.V(4).Infof("Starting grpc connecting...")

		metaInfo, err := c.lister.Get(publicv1alpha1.MetaNameOfKruiseManager)
		if err != nil {
			if errors.IsNotFound(err) {
				klog.Warningf("Not found MetaInfo %s, waiting...", publicv1alpha1.MetaNameOfKruiseManager)
			} else {
				klog.Warningf("Failed to get MetaInfo %s: %v, waiting...", publicv1alpha1.MetaNameOfKruiseManager, err)
			}
			continue
		}

		if metaInfo.Status.Ports == nil || metaInfo.Status.Ports.GrpcLeaderElectionPort == 0 {
			klog.Warningf("No grpc port in MetaInfo %s, waiting...", util.DumpJSON(metaInfo))
			continue
		}

		var leader *publicv1alpha1.MetaInfoEndpoint
		for i := range metaInfo.Status.Endpoints {
			e := &metaInfo.Status.Endpoints[i]
			if e.Leader {
				leader = e
				break
			}
		}
		if leader == nil {
			klog.Warningf("No leader in MetaInfo %s, waiting...", util.DumpJSON(metaInfo))
			continue
		}

		addr := fmt.Sprintf("%s:%d", leader.PodIP, metaInfo.Status.Ports.GrpcLeaderElectionPort)
		func() {
			var opts []grpc.DialOption
			opts = append(opts, grpc.WithInsecure())
			grpcConn, err := grpc.Dial(addr, opts...)
			if err != nil {
				klog.Errorf("Failed to grpc connect to kruise-manager %s addr %s: %v", leader.Name, addr, err)
				return
			}
			ctx, cancel := context.WithCancel(parentCtx)
			defer func() {
				cancel()
				_ = grpcConn.Close()
			}()

			grpcCtrlMeshClient := ctrlmeshproto.NewControllerMeshClient(grpcConn)
			connStream, err := grpcCtrlMeshClient.Register(ctx)
			if err != nil {
				klog.Errorf("Failed to register to kruise-manager %s addr %s: %v", leader.Name, addr, err)
				return
			}

			if err := c.syncing(connStream, initChan); err != nil {
				klog.Errorf("Failed syncing grpc connection to kruise-manager %s addr %s: %v", leader.Name, addr, err)
			}
		}()
	}
}

func (c *GrpcClient) syncing(connStream ctrlmeshproto.ControllerMesh_RegisterClient, initChan chan struct{}) error {
	// firstly
	status := &ctrlmeshproto.Status{SelfInfo: selfInfo}
	if err := connStream.Send(status); err != nil {
		return fmt.Errorf("send first status %s error: %v", util.DumpJSON(status), err)
	}

	if _, err := c.recv(connStream); err != nil {
		return err
	}
	onceInit.Do(func() { close(initChan) })

	for {
		status := &ctrlmeshproto.Status{SpecHash: c.specHash}
		if err := connStream.Send(status); err != nil {
			return fmt.Errorf("send status %s error: %v", util.DumpJSON(status), err)
		}

		for {
			updated, err := c.recv(connStream)
			if err != nil {
				return err
			}
			if updated {
				break
			}
		}
	}
}

func (c *GrpcClient) recv(connStream ctrlmeshproto.ControllerMesh_RegisterClient) (bool, error) {
	spec, err := connStream.Recv()
	if err != nil {
		return false, fmt.Errorf("receive spec error: %v", err)
	}

	specHash := ctrlmeshutil.CalculateHashForProtoSpec(spec)
	if reflect.DeepEqual(specHash, c.specHash) {
		return false, nil
	}

	route := &ctrlmeshproto.InternalRoute{}
	if err = route.DecodeFrom(spec.Route); err != nil {
		klog.Errorf("Failed to decode from spec route %v: %v", util.DumpJSON(spec.Route), err)
		return false, nil
	}

	c.mu.Lock()
	c.route = route
	c.endpoints = spec.Endpoints
	c.specHash = specHash
	c.mu.Unlock()

	klog.V(1).Infof("Refresh new proto spec, subset: '%v', globalLimits: %s, subRules: %s, subsets: %s, sensSubsetNamespaces: %+v. endpoints: %s. SpecHash from %+v to %+v",
		route.Subset, util.DumpJSON(route.GlobalLimits), util.DumpJSON(route.SubRules), util.DumpJSON(route.Subsets), route.SensSubsetNamespaces,
		util.DumpJSON(spec.Endpoints), c.specHash, specHash)
	return true, nil
}

func (c *GrpcClient) GetProtoSpec() (*ctrlmeshproto.InternalRoute, []*ctrlmeshproto.Endpoint, *ctrlmeshproto.SpecHash) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.route, c.endpoints, c.specHash
}
