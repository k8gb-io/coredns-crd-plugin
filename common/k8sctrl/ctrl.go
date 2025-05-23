package k8sctrl

/*
Copyright 2022 The k8gb Contributors

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

import (
	"context"
	"net"
	"strings"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	clog "github.com/coredns/coredns/plugin/pkg/log"
	dnsendpoint "github.com/k8gb-io/coredns-crd-plugin/extdns"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/external-dns/endpoint"
)

type KubeController struct {
	client      dnsendpoint.ExtDNSInterface
	controllers []cache.SharedIndexInformer
	labelFilter string
	hasSynced   bool
	epc         cache.SharedIndexInformer
}

type LookupEndpoint func(indexKey string, clientIP net.IP, geoDataFilePath string, geoDataFieldPath ...string) (result LocalDNSEndpoint)

type ResourceWithLookup struct {
	Name   string
	Lookup LookupEndpoint
}

const (
	defaultResyncPeriod   = 0
	endpointHostnameIndex = "endpointHostname"
)

var log = clog.NewWithPlugin("k8s controller")

var Resources = struct {
	DNSEndpoint *ResourceWithLookup
}{
	DNSEndpoint: &ResourceWithLookup{
		Name: "DNSEndpoint",
	},
}

func NewKubeController(ctx context.Context, c *dnsendpoint.ExtDNSClient, label string) *KubeController {
	ctrl := &KubeController{
		client:      c,
		labelFilter: label,
	}
	endpointController := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc:  endpointLister(ctx, ctrl.client, core.NamespaceAll, ctrl.labelFilter),
			WatchFunc: endpointWatcher(ctx, ctrl.client, core.NamespaceAll, ctrl.labelFilter),
		},
		&endpoint.DNSEndpoint{},
		defaultResyncPeriod,
		cache.Indexers{endpointHostnameIndex: endpointHostnameIndexFunc},
	)
	ctrl.epc = endpointController
	Resources.DNSEndpoint.Lookup = ctrl.getEndpointByName
	ctrl.controllers = append(ctrl.controllers, endpointController)
	return ctrl
}

func (ctrl *KubeController) Run() {
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

func (ctrl *KubeController) getEndpointByName(host string, clientIP net.IP, geoDataFilePath string, geoDataFieldPath ...string) (lep LocalDNSEndpoint) {
	log.Infof("Index key %+v", host)
	endpoints := ctrl.getEndpointsByCaseInsensitiveName(host, clientIP, geoDataFilePath, geoDataFieldPath...)
	lep = ctrl.margeLocalDNSEndpoints(host, endpoints)
	return lep
}

// The function tries to find all case sensitive variants.  Returns a map where the call is hostname and the value is LocalDNSEndpoint
func (ctrl *KubeController) getEndpointsByCaseInsensitiveName(host string, clientIP net.IP, geoDataFilePath string, geoDataFieldPath ...string) (result map[string]LocalDNSEndpoint) {

	// The function extracts LocalDNSEndpoints from *DNSEndpoint. The function is hardwired with a case-sensitive extraction scenario and is only used in a
	// single location, so it is currently declared inside the calling function.
	extractLocalEndpoints := func(ep *endpoint.DNSEndpoint, ip net.IP, host string, geoDataFieldPath ...string) (result []LocalDNSEndpoint) {
		result = []LocalDNSEndpoint{}
		for _, e := range ep.Spec.Endpoints {
			if strings.EqualFold(e.DNSName, host) {
				r := LocalDNSEndpoint{}
				r.DNSName = e.DNSName
				r.Labels = e.Labels
				r.TTL = e.RecordTTL
				r.Targets = e.Targets
				if e.Labels["strategy"] == "geoip" {
					targets := r.extractGeo(e, ip, geoDataFilePath, geoDataFieldPath...)
					if len(targets) > 0 {
						r.Targets = targets
					}
				}
				result = append(result, r)
			}
		}
		return result
	}

	epList := ctrl.epc.GetIndexer().List()
	result = make(map[string]LocalDNSEndpoint, 0)
	for _, obj := range epList {
		ep := obj.(*endpoint.DNSEndpoint)
		extracts := extractLocalEndpoints(ep, clientIP, host, geoDataFieldPath...)
		for _, extracted := range extracts {
			if strings.EqualFold(extracted.DNSName, host) {
				result[extracted.DNSName] = extracted
				log.Debugf("including DNSEndpoint: %s", extracted.String())
			}
		}
	}
	return result
}

func (ctrl *KubeController) margeLocalDNSEndpoints(host string, endpoints map[string]LocalDNSEndpoint) LocalDNSEndpoint {
	result := LocalDNSEndpoint{
		DNSName: host,
	}
	result.Labels = endpoints[host].Labels
	result.TTL = endpoints[host].TTL
	for _, ep := range endpoints {
		result.Targets = append(result.Targets, ep.Targets...)
	}
	return result
}
