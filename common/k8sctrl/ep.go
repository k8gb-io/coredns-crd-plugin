package k8sctrl

import (
	"net"

	"github.com/oschwald/maxminddb-golang"
	"sigs.k8s.io/external-dns/endpoint"
)

type EndpointResult struct {
	Targets []string
	TTL     endpoint.TTL
	Labels  map[string]string
	DNSName string
}

func newEndpoint(ep *endpoint.Endpoint, ip net.IP, host string) (result EndpointResult) {
	result = EndpointResult{}
	if ep == nil {
		return result
	}
	result.Targets = ep.Targets
	if ep.Labels["strategy"] == "geoip" {
		result.Targets = extractGeo(ep, ip)
	}
	result.Labels = ep.Labels
	result.DNSName = host
	result.TTL = ep.RecordTTL
	return result
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

