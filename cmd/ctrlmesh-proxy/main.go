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

package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/openkruise/kruise/pkg/ctrlmesh/proxy/grpcclient"

	"github.com/openkruise/kruise/pkg/client"
	"github.com/openkruise/kruise/pkg/ctrlmesh/constants"
	apiserverproxy "github.com/openkruise/kruise/pkg/ctrlmesh/proxy/apiserver"
	"github.com/openkruise/kruise/pkg/ctrlmesh/proxy/metrics"
	webhookproxy "github.com/openkruise/kruise/pkg/ctrlmesh/proxy/webhook"
	"github.com/openkruise/kruise/pkg/util"
	"github.com/openkruise/kruise/pkg/util/healthz"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

var (
	metricsHealthPort  = flag.Int(constants.ProxyMetricsHealthPortFlag, constants.ProxyMetricsHealthPort, "Port to bind 0.0.0.0 and serve metric endpoint/healthz/pprof.")
	proxyApiserverPort = flag.Int(constants.ProxyApiserverPortFlag, constants.ProxyApiserverPort, "Port to bind localhost and proxy the requests to apiserver.")
	proxyWebhookPort   = flag.Int(constants.ProxyWebhookPortFlag, constants.ProxyWebhookPort, "Port to bind 0.0.0.0 and proxy the requests to webhook.")

	leaderElectionName = flag.String(constants.ProxyLeaderElectionNameFlag, "", "The name of leader election.")
	webhookServePort   = flag.Int(constants.ProxyWebhookServePortFlag, 0, "Port that the real webhook binds, 0 means no proxy for webhook.")
	webhookCertDir     = flag.String(constants.ProxyWebhookCertDirFlag, "", "The directory where the webhook certs generated or mounted.")

	cfg           *rest.Config
	grpcClient    = grpcclient.New()
	stop          = signals.SetupSignalHandler()
	healthzProber = healthz.NewHealthz()
)

func main() {
	flag.Parse()

	if os.Getenv(constants.EnvPodNamespace) == "" || os.Getenv(constants.EnvPodName) == "" {
		klog.Fatalf("Environment %s=%s %s=%s not exist.",
			constants.EnvPodNamespace, os.Getenv(constants.EnvPodNamespace), constants.EnvPodName, os.Getenv(constants.EnvPodName))
	}

	var err error
	if err = setupConfig(); err != nil {
		klog.Fatalf("Failed to setup config: %v", err)
	}
	if err = client.NewRegistry(cfg); err != nil {
		klog.Fatalf("Failed to new client for config %s: %v", util.DumpJSON(cfg), err)
	}

	if err = grpcClient.Start(stop); err != nil {
		klog.Fatalf("Failed to start grpcclient: %v", err)
	}

	stoppedWebhook := startWebhookProxy()
	stoppedApiserver := startApiserverProxy()

	serveHTTP()
	if stoppedWebhook != nil {
		select {
		case <-stoppedWebhook:
			klog.Infof("Webhook proxy stopped")
		}
	}
	select {
	case <-stoppedApiserver:
		klog.Infof("Apiserver proxy stopped")
	}
}

func startWebhookProxy() <-chan struct{} {
	if *webhookServePort <= 0 {
		klog.Infof("Skip proxy webhook for webhook serve port not set")
		return nil
	}

	opts := &webhookproxy.Options{
		CertDir:     *webhookCertDir,
		BindPort:    *proxyWebhookPort,
		WebhookPort: *webhookServePort,
		GrpcClient:  grpcClient,
	}
	proxy := webhookproxy.NewProxy(opts)
	healthzProber.RegisterFunc("webhookProxy", proxy.HealthFunc)

	stopped, err := proxy.Start(stop)
	if err != nil {
		klog.Fatalf("Failed to start webhook proxy: %v", err)
	}
	return stopped
}

func startApiserverProxy() <-chan struct{} {
	opts := apiserverproxy.NewOptions()
	opts.Config = rest.CopyConfig(cfg)
	opts.SecureServingOptions.ServerCert.CertKey.KeyFile = "/var/run/secrets/kubernetes.io/serviceaccount/ctrlmesh/tls.key"
	opts.SecureServingOptions.ServerCert.CertKey.CertFile = "/var/run/secrets/kubernetes.io/serviceaccount/ctrlmesh/tls.crt"
	opts.SecureServingOptions.BindAddress = net.ParseIP("127.0.0.1")
	opts.SecureServingOptions.BindPort = *proxyApiserverPort
	opts.LeaderElectionName = *leaderElectionName
	opts.GrpcClient = grpcClient
	errs := opts.Validate()
	if len(errs) > 0 {
		klog.Fatalf("Failed to validate apiserver-proxy options %s: %v", util.DumpJSON(opts), errs)
	}

	proxy, err := apiserverproxy.NewProxy(opts)
	if err != nil {
		klog.Fatalf("Failed to new apiserver proxy: %v", err)
	}

	stopped, err := proxy.Start(stop)
	if err != nil {
		klog.Fatalf("Failed to start apiserver proxy: %v", err)
	}
	return stopped
}

func serveHTTP() {
	mux := http.DefaultServeMux
	mux.Handle("/metrics", promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{
		ErrorHandling: promhttp.HTTPErrorOnError,
	}))
	mux.HandleFunc("/healthz", healthzProber.Handler)

	server := http.Server{
		Handler: mux,
	}

	// Run the server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *metricsHealthPort))
	if err != nil {
		klog.Fatalf("Failed to listen on :%d: %v", *metricsHealthPort, err)
	}
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			klog.Fatalf("Failed to serve HTTP on :%d: %v", *metricsHealthPort, err)
		}
	}()

	// Shutdown the server when stop is closed
	<-stop
	if err := server.Shutdown(context.Background()); err != nil {
		klog.Fatalf("Serve HTTP shutting down on :%d: %v", *metricsHealthPort, err)
	}
}

func setupConfig() error {
	const (
		tokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
		//rootCAFile = "/var/run/secrets/kubernetes.io/serviceaccount/..data/ca.crt"
	)
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		return rest.ErrNotInCluster
	}

	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return err
	}

	tlsClientConfig := rest.TLSClientConfig{Insecure: true}

	//if _, err := certutil.NewPool(rootCAFile); err != nil {
	//	klog.Errorf("Expected to load root CA config from %s, but got err: %v", rootCAFile, err)
	//} else {
	//	tlsClientConfig.CAFile = rootCAFile
	//}

	cfg = &rest.Config{
		// TODO: switch to using cluster DNS.
		Host:            "https://" + net.JoinHostPort(host, port),
		TLSClientConfig: tlsClientConfig,
		BearerToken:     string(token),
		BearerTokenFile: tokenFile,

		Burst:   3000,
		QPS:     2000.0,
		Timeout: time.Duration(300) * time.Second,
	}
	klog.V(3).Infof("Starting with rest config: %v", util.DumpJSON(cfg))

	return nil
}
