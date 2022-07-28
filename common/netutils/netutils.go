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
