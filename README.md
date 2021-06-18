# k8s_crd

A CoreDNS plugin that is very similar to [k8s_external](https://coredns.io/plugins/k8s_external/) but supporting DNSEndpoint external resource.

**This project is a modification of [k8s_gateway](https://github.com/ori-edge/k9s_gateway) plugin, adopted with DNSEndpoint client.**

This plugin relies on it's own connection to the k8s API server and doesn't share any code with the existing [kubernetes](https://coredns.io/plugins/kubernetes/) plugin. The assumption is that this plugin can now be deployed as a separate instance (alongside the internal kube-dns) and act as a single external DNS interface into your Kubernetes cluster(s).


## Description

`k8s_crd` resolves Kubernetes resources with their external IP addresses based on zones specified in the configuration. This plugin will resolve the following type of resources:

| Kind | Matching Against | External IPs are from |
| ---- | ---------------- | -------- |
| DNSEndponit | all FQDNs from `spec.endpoints.dnszone` matching configured zones | `.spec.endpoints.dnszone.targets` |


Currently only supports A-type queries, all other queries result in NODATA responses.

This plugin is **NOT** supposed to be used for intra-cluster DNS resolution and does not contain the default upstream [kubernetes](https://coredns.io/plugins/kubernetes/) plugin.

## Install

The recommended installation method is using the helm chart provided in the repo:

```
helm install exdns ./charts/coredns
```
## Configure

```
k8s_crd [ZONE...]
```

Optionally, you can specify what kind of resources to watch, default TTL to return in response and a default name to use for zone apex, e.g.

```
k8s_crd example.com {
    resources DNSEndpoint
    ttl 10
    apex dns1
}
```

## Resolving order
In case dnsendpoint object's target have label `strategy: geoip` CoreDNS k8s_crd plugin will respond in a special way:
* assuming record have multiple IPs associated with it, and dns Message comes with edns0 CLIENT-SUBNET option.
* CoreDNS will compare "DC" tag for IP exctracted from CLIENT-SUBNET option against available endpoint.targets
* Return only IPs where tags match
* If IP have no common tag, all entries are returned.
* CoreDNS must be supplied with specially crafted GeoIP database in MaxMind DB format and mounted as `/geoip.mmdb` Refer to `terratest/geo` for examples.

## Build

### With compile-time configuration file

```
$ git clone https://github.com/coredns/coredns
$ cd coredns
$ vim plugin.cfg
# Replace lines with kubernetes and k8s_external with k8s_crd:github.com/absaoss/k8s_crd
$ go generate
$ go build
$ ./coredns -plugins | grep k8s_crd
```

### With external golang source code
```
$ git clone https://github.com/absaoss/k8s_crd.git
$ cd k8s_crd
$ go build cmd/coredns.go
$ ./coredns -plugins | grep k8s_crd
```

For more details refer to [this CoreDNS doc](https://coredns.io/2017/07/25/compile-time-enabling-or-disabling-plugins/)


## Notes regarding Zone Apex and NS server resolution

Due to the fact that there is not nice way to discover NS server's own IP to respond to A queries, as a wokaround, it's possible to pass the name of the LoadBalancer service used to expose the CoreDNS instance as an environment variable `EXTERNAL_SVC`. If not set, the default fallback value of `external-dns.kube-system` will be used to look up the external IP of the CoreDNS service.
