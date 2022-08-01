package wrr

import (
	"context"
	"fmt"
	"net"

	"github.com/AbsaOSS/k8s_crd/common/k8sctrl"
	"github.com/AbsaOSS/k8s_crd/common/netutils"
	"github.com/coredns/coredns/request"

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

	a, _, _ := parseAnswerSection(r.Answer)

	var ep = k8sctrl.Resources.DNSEndpoint.Lookup(indexKey, clientIP)

	fmt.Println(a, ep.Labels)
	return dns.RcodeSuccess, nil
}

func (wrr *WeightRoundRobin) Name() string { return thisPlugin }
