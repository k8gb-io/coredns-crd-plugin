/*
Copyright 2021 ABSA Group Limited

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/
package gateway

import (
	"context"
	"net"
	"os"
	"strings"

	"k8s.io/client-go/tools/clientcmd"

	// "os"

	dnsendpoint "github.com/AbsaOSS/k8s_crd/extdns"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	endpoint "sigs.k8s.io/external-dns/endpoint"
	// "k8s.io/client-go/tools/clientcmd"
)

const (
	defaultResyncPeriod   = 0
	endpointHostnameIndex = "endpointHostname"
)

type geo struct {
	DC string `maxminddb:"datacenter"`
}

// KubeController stores the current runtime configuration and cache
type KubeController struct {
	client      dnsendpoint.ExtDNSInterface
	controllers []cache.SharedIndexInformer
	labelFilter string
	hasSynced   bool
}

func newKubeController(ctx context.Context, c *dnsendpoint.ExtDNSClient, label string) *KubeController {

	log.Infof("Starting k8s_crd controller")

	ctrl := &KubeController{
		client:      c,
		labelFilter: label,
	}
	if resource := lookupResource("DNSEndpoint"); resource != nil {
		endpointController := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc:  endpointLister(ctx, ctrl.client, core.NamespaceAll, label),
				WatchFunc: endpointWatcher(ctx, ctrl.client, core.NamespaceAll, label),
			},
			&endpoint.DNSEndpoint{},
			defaultResyncPeriod,
			cache.Indexers{endpointHostnameIndex: endpointHostnameIndexFunc},
		)
		resource.lookup = lookupEndpointIndex(endpointController)
		ctrl.controllers = append(ctrl.controllers, endpointController)
	}

	return ctrl
}

func (ctrl *KubeController) run() {
	stopCh := make(chan struct{})
	defer close(stopCh)

	var synced []cache.InformerSynced

	for _, ctrl := range ctrl.controllers {
		go ctrl.Run(stopCh)
		synced = append(synced, ctrl.HasSynced)
	}

	if !cache.WaitForCacheSync(stopCh, synced...) {
		ctrl.hasSynced = false
	}
	log.Infof("Synced all required resources")
	ctrl.hasSynced = true

	<-stopCh
}

// HasSynced returns true if all controllers have been synced
func (ctrl *KubeController) HasSynced() bool {
	return ctrl.hasSynced
}

// RunKubeController kicks off the k8s controllers
func RunKubeController(ctx context.Context, c *Gateway) (*KubeController, error) {
	// config, err := rest.InClusterConfig()

	// Helpful to run coredns locally
	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		return nil, err
	}

	err = dnsendpoint.AddToScheme(scheme.Scheme)
	if err != nil {
		return nil, err
	}

	kubeClient, err := dnsendpoint.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	ctrl := newKubeController(ctx, kubeClient, c.Filter)

	go ctrl.run()

	return ctrl, nil

}

func endpointLister(ctx context.Context, c dnsendpoint.ExtDNSInterface, ns, label string) func(meta.ListOptions) (runtime.Object, error) {
	return func(opts meta.ListOptions) (runtime.Object, error) {
		opts.LabelSelector = label
		return c.DNSEndpoints(ns).List(ctx, opts)
	}
}

func endpointWatcher(ctx context.Context, c dnsendpoint.ExtDNSInterface, ns, label string) func(meta.ListOptions) (watch.Interface, error) {
	return func(opts meta.ListOptions) (watch.Interface, error) {
		opts.LabelSelector = label
		return c.DNSEndpoints(ns).Watch(ctx, opts)
	}
}

func endpointHostnameIndexFunc(obj interface{}) ([]string, error) {
	ep, ok := obj.(*endpoint.DNSEndpoint)
	if !ok {
		return []string{}, nil
	}

	var hostnames []string
	for _, rule := range ep.Spec.Endpoints {
		log.Infof("Adding index %s for endpoints %s", rule.DNSName, ep.Name)
		hostnames = append(hostnames, rule.DNSName)
	}
	return hostnames, nil
}

func lookupEndpointIndex(ctrl cache.SharedIndexInformer) func(string, net.IP) ([]string, endpoint.TTL) {
	return func(indexKey string, clientIP net.IP) (result []string, ttl endpoint.TTL) {
		log.Infof("Index key %+v", indexKey)
		objs, _ := ctrl.GetIndexer().ByIndex(endpointHostnameIndex, strings.ToLower(indexKey))
		for _, obj := range objs {
			endpoint := obj.(*endpoint.DNSEndpoint)
			result, ttl = fetchEndpointTargets(endpoint.Spec.Endpoints, indexKey, clientIP)
		}
		return
	}
}
