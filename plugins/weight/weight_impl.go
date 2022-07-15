package weight

import "github.com/miekg/dns"

type weight struct{}

func newWeight() *weight {
	return &weight{}
}

func (w *weight) update() ([]dns.RR, error) {
	return nil, nil
}
