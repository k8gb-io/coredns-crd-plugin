package k8sctrl

import (
	"net"

	"github.com/oschwald/maxminddb-golang"
	"sigs.k8s.io/external-dns/endpoint"
)

type geo struct {
	DC string `maxminddb:"datacenter"`
}

type LocalDNSEndpoint struct {
	Targets []string
	TTL     endpoint.TTL
	Labels  map[string]string
	DNSName string
}

func extractLocalEndpoint(ep *endpoint.DNSEndpoint, ip net.IP, host string) (result LocalDNSEndpoint) {
	result = LocalDNSEndpoint{}
	for _, e := range ep.Spec.Endpoints {
		if e.DNSName == host {
			result.DNSName = host
			result.Labels = e.Labels
			result.TTL = e.RecordTTL
			result.Targets = e.Targets
			if e.Labels["strategy"] == "geoip" {
				targets := result.extractGeo(e, ip)
				if len(targets) > 0 {
					result.Targets = targets
				}
			}
			break
		}
	}
	return result
}

func (lep LocalDNSEndpoint) isEmpty() bool {
	return len(lep.Targets) == 0 && (len(lep.Labels) == 0) && (lep.TTL == 0)
}

func (lep LocalDNSEndpoint) extractGeo(endpoint *endpoint.Endpoint, clientIP net.IP) (result []string) {
	db, err := maxminddb.Open("geoip.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	clientGeo := &geo{}
	err = db.Lookup(clientIP, clientGeo)
	if err != nil {
		return nil
	}

	if clientGeo.DC == "" {
		log.Infof("empty DC %+v", clientGeo)
		return result
	}

	log.Infof("clientDC: %+v", clientGeo)

	for _, ip := range endpoint.Targets {
		geoData := &geo{}
		log.Infof("processing IP %+v", ip)
		err = db.Lookup(net.ParseIP(ip), geoData)
		if err != nil {
			log.Error(err)
			continue
		}

		log.Infof("IP info: %+v", geoData.DC)
		if clientGeo.DC == geoData.DC {
			result = append(result, ip)
		}
	}
	return result
}
