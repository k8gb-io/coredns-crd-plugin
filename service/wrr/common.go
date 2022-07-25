package wrr

import "github.com/miekg/dns"

// parseAnswerSection converts []dns.RR into map of A or AAAA records and slice containing all except A or AAAA
func parseAnswerSection(arr []dns.RR) (ipmap map[string]dns.RR, ip []string, noip []dns.RR) {
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
