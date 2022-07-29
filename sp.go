package k8scrd

import (
	"github.com/AbsaOSS/k8s_crd/service/gateway"
	"github.com/AbsaOSS/k8s_crd/service/wrr"
)

type args struct {
	annotation     string
	apex           string
	filter         string
	kubecontroller string
	negttl         uint32
	resources      []string
	ttl            uint32
	zones          []string
}

func (a args) provideGatewayService() (gw *gateway.Gateway) {
	opts := gateway.NewGatewayOpts(a.annotation, a.apex, a.ttl, a.negttl, a.resources, a.zones)
	gw = gateway.NewGateway(opts)
	return gw
}

func (a args) provideWrrService() (w *wrr.WeightRoundRobin) {
	w = wrr.NewWeightRoundRobin()
	return w
}
