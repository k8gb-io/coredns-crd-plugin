package weight

import (
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Weight struct {
	w *weight
}

func NewWeight() *Weight {
	return &Weight{
		newWeight(),
	}
}

func (s *Weight) Shuffle(req request.Request, res *dns.Msg) ([]dns.RR, error) {
	return s.w.update()
}
