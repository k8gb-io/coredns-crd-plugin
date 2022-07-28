package k8sctrl

import (
	"net"
	"strings"

	"sigs.k8s.io/external-dns/endpoint"
)

func (ctrl *KubeController) lookupEndpointIndex(indexKey string, clientIP net.IP) (result []string, ttl endpoint.TTL) {
	log.Infof("Index key %+v", indexKey)
	objs, _ := ctrl.epc.GetIndexer().ByIndex(endpointHostnameIndex, strings.ToLower(indexKey))
	for _, obj := range objs {
		endpoint := obj.(*endpoint.DNSEndpoint)
		result, ttl = fetchEndpointTargets(endpoint.Spec.Endpoints, indexKey, clientIP)
	}
	return
}

func fetchEndpointTargets(endpoints []*endpoint.Endpoint, host string, ip net.IP) (results []string, ttl endpoint.TTL) {
	for _, ep := range endpoints {
		if ep.DNSName == host {
			ttl = ep.RecordTTL
			log.Info("oldone: ",ep.DNSName," LABELS: ", ep.Labels)
			if ep.Labels["strategy"] == "geoip" {
				log.Info("oldone: GEO")
				results = extractGeo(ep, ip)
				log.Info("oldone:",results," ",ip.String())
				if len(results) > 0 {
					return
				}
			} else {
				log.Info("oldone: NOGEO")
			}
			results = ep.Targets
		}
	}
	return
}
