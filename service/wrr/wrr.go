package wrr

import (
	"context"
	"fmt"

	"github.com/miekg/dns"
)

type WeightRoundRobin struct {
}

const thisPlugin = "wrr"

func NewWeightRoundRobin() *WeightRoundRobin {
	return &WeightRoundRobin{}
}

func (wrr *WeightRoundRobin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	a, _, _ := parseAnswerSection(r.Answer)
	fmt.Println(a)
	return dns.RcodeSuccess, nil
}

func (wrr *WeightRoundRobin) Name() string { return thisPlugin }
