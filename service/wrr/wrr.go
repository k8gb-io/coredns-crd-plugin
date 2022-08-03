package wrr

import (
	"context"
	"fmt"
	"net"

	"github.com/AbsaOSS/k8s_crd/common/k8sctrl"
	"github.com/AbsaOSS/k8s_crd/common/netutils"

	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/k8gb-io/go-weight-shuffling/gows"
	"github.com/miekg/dns"
)

type WeightRoundRobin struct {
}

const thisPlugin = "wrr"

var log = clog.NewWithPlugin(thisPlugin)

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
	if !g.hasWeights() {
		return dns.RcodeSuccess, nil
	}

	ws, err := gows.NewWS(g.pdf())
	if err != nil {
		err = fmt.Errorf("error create distribution (%s)", err)
		return dns.RcodeServerFailure, err
	}

	vector := ws.PickVector()
	g.shuffle(vector)

	m := new(dns.Msg)
	m.SetReply(state.Req)
	m.Answer = wrr.updateAnswers(g, r.Answer)
	if err := w.WriteMsg(m); err != nil {
		log.Errorf("Failed to send a response: %s", err)
	}
	return dns.RcodeSuccess, err
}

func (wrr *WeightRoundRobin) Name() string { return thisPlugin }

// updateAnswers set order of answers based on groups. The function doesn't handle
// the fact that answers does not match the weight-labels in the endpoint because
// other services can add or remove answers.
func (wrr *WeightRoundRobin) updateAnswers(g groups, answers []dns.RR) (newAnswers []dns.RR) {
	order := g.asSlice()
	targets, _, noa := netutils.ParseAnswerSection(answers)
	newAnswers = append(newAnswers, noa...)
	for _, ip := range order {
		if rr, found := targets[ip]; found {
			newAnswers = append(newAnswers, rr)
			continue
		}
		log.Infof("[%s] exist as target but missing in labels", ip)
	}
	return newAnswers
}
