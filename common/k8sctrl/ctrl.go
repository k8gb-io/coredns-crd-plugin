package k8sctrl

import (
	"context"
	"net"
	"strings"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	dnsendpoint "github.com/AbsaOSS/k8s_crd/extdns"
	clog "github.com/coredns/coredns/plugin/pkg/log"
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
	resources   []*ResourceWithIndex
	epc         cache.SharedIndexInformer
}

type LookupEndpoint func(indexKey string, clientIP net.IP) (result LocalDNSEndpoint)

type ResourceWithIndex struct {
	Name   string
	Lookup LookupEndpoint
}

const (
	defaultResyncPeriod   = 0
	endpointHostnameIndex = "endpointHostname"
)

// TODO: is new logger instance necessary
var log = clog.NewWithPlugin("k8s controller")

var Resources = struct {
	DNSEndpoint *ResourceWithIndex
}{
	DNSEndpoint: &ResourceWithIndex{
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
			ListFunc:  endpointLister(ctx, ctrl.client, core.NamespaceAll, label),
			WatchFunc: endpointWatcher(ctx, ctrl.client, core.NamespaceAll, label),
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

func (ctrl *KubeController) getEndpointByName(host string, clientIP net.IP) (lep LocalDNSEndpoint) {
	log.Infof("Index key %+v", host)
	objs, _ := ctrl.epc.GetIndexer().ByIndex(endpointHostnameIndex, strings.ToLower(host))
	for _, obj := range objs {
		ep := obj.(*endpoint.DNSEndpoint)
		lep = extractLocalEndpoint(ep, clientIP, host)
		if !lep.isEmpty() {
			break
		}
	}
	return lep
}
