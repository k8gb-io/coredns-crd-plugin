package wrr

import (
	"context"
	"fmt"
	"net"

	"github.com/AbsaOSS/k8s_crd/common/k8sctrl"
	"github.com/AbsaOSS/k8s_crd/common/netutils"

	"github.com/coredns/coredns/request"
	"github.com/k8gb-io/go-weight-shuffling/gows"
	"github.com/miekg/dns"
)

type WeightRoundRobin struct {
}

const thisPlugin = "wrr"

func NewWeightRoundRobin() *WeightRoundRobin {
	return &WeightRoundRobin{}
}

func (wrr *WeightRoundRobin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	var clientIP net.IP
	state := request.Request{W: w, Req: r}
	clientIP = netutils.ExtractEdnsSubnet(r)
	indexKey := netutils.StripClosingDot(state.QName())

	var ep = k8sctrl.Resources.DNSEndpoint.Lookup(indexKey, clientIP)
	g, err := parseGroups(ep.Labels)
	if err != nil {
		err = fmt.Errorf("error parsing lables (%s)", err)
		return dns.RcodeServerFailure, err
	}

	ws, err := gows.NewWS(g.pdf())
	if err != nil {
		err = fmt.Errorf("error create distribution (%s)", err)
		return dns.RcodeServerFailure, err
	}

	vector := ws.PickVector()

	g.shuffle(vector)

	return dns.RcodeSuccess, nil
}

func (wrr *WeightRoundRobin) Name() string { return thisPlugin }

// strategy:roundRobin
// weight-eu-0-50:172.18.0.5
// weight-eu-1-50:172.18.0.6
// weight-us-0-50:172.18.0.8
// weight-us-1-50:172.18.0.9
