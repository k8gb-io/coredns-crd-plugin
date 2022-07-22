package gateway

import (
	"net"

	"github.com/oschwald/maxminddb-golang"
	endpoint "sigs.k8s.io/external-dns/endpoint"
)

func fetchEndpointTargets(endpoints []*endpoint.Endpoint, host string, ip net.IP) (results []string, ttl endpoint.TTL) {
	for _, ep := range endpoints {
		if ep.DNSName == host {
			ttl = ep.RecordTTL
			if ep.Labels["strategy"] == "geoip" {
				results = extractGeo(ep, ip)
				if len(results) > 0 {
					return
				}
			}
			results = ep.Targets
		}
	}
	return
}

func extractGeo(endpoint *endpoint.Endpoint, clientIP net.IP) (result []string) {
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
