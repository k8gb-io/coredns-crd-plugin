package netutils

import (
	"net"
	"strings"

	"github.com/miekg/dns"
)

func ExtractEdnsSubnet(msg *dns.Msg) net.IP {
	edns := msg.IsEdns0()
	if edns == nil {
		return nil
	}
	for _, o := range edns.Option {
		if o.Option() == dns.EDNS0SUBNET {
			subnet := o.(*dns.EDNS0_SUBNET)
			return subnet.Address
		}
	}
	return nil
}

// StripClosingDot strips the closing dot unless it's "."
func StripClosingDot(s string) string {
	if len(s) > 1 {
		return strings.TrimSuffix(s, ".")
	}
	return s
}

func TargetToIP(targets []string) (ips []net.IP) {
	for _, ip := range targets {
		ips = append(ips, net.ParseIP(ip))
	}
	return
}

// ParseAnswerSection converts []dns.RR into map of A or AAAA records and slice containing all except A or AAAA
func ParseAnswerSection(arr []dns.RR) (ipmap map[string]dns.RR, ip []string, noip []dns.RR) {
	ipmap = make(map[string]dns.RR)
	ip = make([]string, 0)
	noip = make([]dns.RR, 0)
	for _, r := range arr {
		switch r.Header().Rrtype {
		case dns.TypeA:
			a := r.(*dns.A).A.String()
			ipmap[a] = r
			ip = append(ip, a)
		case dns.TypeAAAA:
			aaaa := r.(*dns.AAAA).AAAA.String()
			ipmap[aaaa] = r
			ip = append(ip, aaaa)
		default:
			noip = append(noip, r)
		}
	}
	return
}
