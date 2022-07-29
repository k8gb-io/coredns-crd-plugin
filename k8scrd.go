package k8scrd

import (
	"context"
	"fmt"

	"github.com/AbsaOSS/k8s_crd/common/k8sctrl"

	"github.com/AbsaOSS/k8s_crd/service"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type K8sCRD struct {
	Next       plugin.Handler
	Filter     string
	Controller *k8sctrl.KubeController
	container  service.PluginContainer
}

func NewK8sCRD() *K8sCRD {
	return &K8sCRD{
		container: service.NewCommonContainer(),
	}
}

// ServeDNS implements the plugin.Handle interface.
func (p *K8sCRD) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	if !p.Controller.HasSynced() {
		// TODO maybe there's a better way to do this? e.g. return an error back to the client?
		return dns.RcodeServerFailure, plugin.Error(thisPlugin, fmt.Errorf("could not sync required resources"))
	}

	err := p.container.Execute(ctx, w, r)
	if err != nil {
		return dns.RcodeServerFailure, err
	}
	return dns.RcodeSuccess, nil
}

func (p *K8sCRD) Name() string {
	return thisPlugin
}
