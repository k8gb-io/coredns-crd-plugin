package k8sctrl

import "sigs.k8s.io/external-dns/endpoint"

type Endpoint struct {
	Targets []string
	TTL     endpoint.TTL
	Labels  map[string]string
	DNSName string
}

type EndpointResult struct {
	Endpoints []Endpoint
}

func (r *EndpointResult) Append(targets []string, labels map[string]string, name string, ttl endpoint.TTL) {
	ep := Endpoint{
		Targets: targets,
		Labels:  labels,
		DNSName: name,
		TTL:     ttl,
	}
	r.Endpoints = append(r.Endpoints, ep)
}

func (r *EndpointResult) IsEmpty() bool {
	return len(r.Endpoints) == 0
}
