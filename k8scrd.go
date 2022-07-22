package k8scrd

import (
	"context"

	"github.com/AbsaOSS/k8s_crd/service"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type K8sCRD struct {
	Next      plugin.Handler
	container service.ServiceContainer
}

func NewK8sCRD() *K8sCRD {
	return &K8sCRD{
		container: service.NewCommonContainer(),
	}
}

// ServeDNS implements the plugin.Handle interface.
func (p *K8sCRD) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	err := p.container.Execute(ctx, w, r)
	if err != nil {
		return dns.RcodeServerFailure, err
	}
	return dns.RcodeSuccess, nil
}

func (p *K8sCRD) Name() string {
	return thisPlugin
}
