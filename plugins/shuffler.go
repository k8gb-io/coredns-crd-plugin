package plugins

import (
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Shuffler interface {
	// Shuffle runs round-robin algorithm.
	// stateless contains incoming request while *msg is response modified by other plugins
	Shuffle(req request.Request, msg *dns.Msg) ([]dns.RR, error)
}
