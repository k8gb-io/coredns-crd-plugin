package wrr

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
	"fmt"
	"net"
	"strings"

	"github.com/k8gb-io/coredns-crd-plugin/common/k8sctrl"
	"github.com/k8gb-io/coredns-crd-plugin/common/netutils"

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

func (wrr *WeightRoundRobin) ServeDNS(_ context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// skipping if no answers
	if len(r.Answer) == 0 {
		log.Debugf("no answers in response; skipping")
		return dns.RcodeSuccess, nil
	}
	var clientIP net.IP
	state := request.Request{W: w, Req: r}
	clientIP = netutils.ExtractEdnsSubnet(r)
	indexKey := netutils.StripClosingDot(state.QName())
	var ep = k8sctrl.Resources.DNSEndpoint.Lookup(indexKey, clientIP, "")
	// weights are not defined, labels doesnt exists
	if len(ep.Labels) == 1 && strings.ToUpper(ep.Labels["strategy"]) == "ROUNDROBIN" {
		roundRobinShuffle(r.Answer)
		if err := w.WriteMsg(r); err != nil {
			return dns.RcodeServerFailure, fmt.Errorf("[random] %s", err)
		}
		_, ansIP, _ := netutils.ParseAnswerSection(r.Answer)
		log.Infof("[random]%s", ansIP)
		return dns.RcodeSuccess, nil
	}
	g, err := parseGroups(ep.Labels)
	if err != nil {
		log.Errorf("error parsing labels (%s)", err)
		return dns.RcodeServerFailure, nil
	}
	if !g.hasWeights() {
		return dns.RcodeSuccess, nil
	}
	// protection from incomplete answers (incomplete responses are generated during initialization,
	// or when DNSEndpoint is not properly adjusted)
	_, ansIP, _ := netutils.ParseAnswerSection(r.Answer)
	if !g.equalIPs(ansIP) {
		log.Infof("[skipping]%s", ansIP)
		log.Debugf("DNSEndpoint labels does not match with DNS answer. DNSEndpoint might not be initialised yet. ep: %v; answers: %v", g.asSlice(), ansIP)
		return dns.RcodeSuccess, nil
	}
	ws := gows.New(g.pdf())
	vector := ws.PickVector(gows.IncludeZeroWeights)
	g.shuffle(vector)
	log.Infof("%v%v: %v", vector, g, g.asSlice())
	m := new(dns.Msg)
	m.SetReply(state.Req)
	m.Answer = wrr.updateAnswers(g, r.Answer)
	if err = w.WriteMsg(m); err != nil {
		log.Errorf("Failed to send a response: %s", err)
		return dns.RcodeServerFailure, nil
	}
	return dns.RcodeSuccess, nil
}

func (wrr *WeightRoundRobin) Name() string { return thisPlugin }

// updateAnswers set order of answers based on groups. The function doesn't handle
// the fact that answers does not match the weight-labels in the endpoint because
// other services can add or remove answers. In such case function produces warning.
func (wrr *WeightRoundRobin) updateAnswers(g groups, answers []dns.RR) (newAnswers []dns.RR) {
	labelIPs := g.asSlice()
	targets, _, noa := netutils.ParseAnswerSection(answers)
	newAnswers = append(newAnswers, noa...)
	for _, ip := range labelIPs {
		if rr, found := targets[ip]; found {
			newAnswers = append(newAnswers, rr)
			continue
		}
		// some plugin before wrr removed any IP(s). It exists in labels but missing in answers
		log.Warningf("[%s] exist as label but missing in incoming message", ip)
	}
	return newAnswers
}

// taken from https://github.com/coredns/coredns/blob/master/plugin/loadbalance/loadbalance.go
func roundRobinShuffle(records []dns.RR) {
	switch l := len(records); l {
	case 0, 1:
		break
	case 2:
		if dns.Id()%2 == 0 {
			records[0], records[1] = records[1], records[0]
		}
	default:
		for j := 0; j < l; j++ {
			p := j + (int(dns.Id()) % (l - j))
			if j == p {
				continue
			}
			records[j], records[p] = records[p], records[j]
		}
	}
}
