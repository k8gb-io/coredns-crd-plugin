# k8s_crd

A CoreDNS plugin that is very similar to [k8s_external](https://coredns.io/plugins/k8s_external/) but supporting DNSEndpoint external resource.

**This project is a modification of [k8s_gateway](https://github.com/ori-edge/k8s_gateway) plugin, adopted with DNSEndpoint client.**

This plugin relies on it's own connection to the k8s API server and doesn't share any code with the existing [kubernetes](https://coredns.io/plugins/kubernetes/) plugin. The assumption is that this plugin can now be deployed as a separate instance (alongside the internal kube-dns) and act as a single external DNS interface into your Kubernetes cluster(s).

## Description

`k8s_crd` resolves Kubernetes resources with their external IP addresses based on zones specified in the configuration. This plugin will resolve the following type of resources:

| Kind | Matching Against | External IPs are from |
| ---- | ---------------- | -------- |
| DNSEndpoint | all FQDNs from `spec.endpoints.dnszone` matching configured zones | `.spec.endpoints.dnszone.targets` |

Currently only supports A-type queries, all other queries result in NODATA responses.

This plugin is **NOT** supposed to be used for intra-cluster DNS resolution and does not contain the default upstream [kubernetes](https://coredns.io/plugins/kubernetes/) plugin.

## Install

The recommended installation method is using the helm chart provided in the repo:

```shell
helm install exdns ./charts/coredns
```

## Configure

```text
k8s_crd [ZONE...]
```

Optionally, you can specify what kind of resources to watch, default TTL to return in response and a default name to use for zone apex, e.g.

```text
k8s_crd example.com {
    ttl 10
    apex dns1
}
```

## Resolving order

### GeoIP

In case dnsEndpoint object's target has a label of `strategy: geoip` CoreDNS `k8s_crd` plugin will respond in a special way:

* Assuming record has multiple IPs associated with it, and DNS message comes with edns0 `CLIENT-SUBNET` option.
* CoreDNS will compare the specified field tag (`datacenter` by default, configured via the `geodatafield` plugin option) for IP extracted from `CLIENT-SUBNET` option against available Endpoint.Targets
* Return only IPs where tags match
* If IP has no common tag, all entries are returned.
* CoreDNS must be supplied with a specially crafted GeoIP database in MaxMind DB format and mounted (at `/geoip.mmdb` by default, configured via the `geodatafilepath` plugin option). Refer to [./terratest/geogen](./terratest/geogen) for examples. Using the MaxMind GeoLite2 database is supported using the necessary `geodatafield` to configure the field to use as required.

The following configuration options are available:

```text
k8s_crd example.com {
    geodatafilepath /geoip.mmdb
    geodatafield country.iso_code
    ...
}
```

### Weight Round Robin

To enable the weight round robin you have to set the configuration to weight load-balancer:

```text
k8s_crd example.com {
    loadbalance weight
    ...
}
```

The dnsEndpoint must also contain information about the percentage distribution per region
and their IP addresses. Thanks to this, the weight round-robin module will know in which
order to return IP addresses. Addresses with high probability will often be at the top of
DNS responses, while those with low probability will be at the bottom.

```yaml
labels:
    strategy: roundrobin
    weight-eu-0-50: 10.0.0.1
    weight-eu-1-50: 10.0.0.2
    weight-za-0-0:  10.10.0.1
    weight-us-0-50: 10.20.0.1
```

For more information about balancing, please visit our [go-weight-shuffling](https://github.com/k8gb-io/go-weight-shuffling
) module.

## Build

### With compile-time configuration file

```shell
git clone https://github.com/coredns/coredns
cd coredns
vim plugin.cfg
# Replace lines with kubernetes and k8s_external with k8s_crd:github.com/absaoss/k8s_crd
go generate
go build
./coredns -plugins | grep k8s_crd
```

### With external golang source code

```shell
git clone https://github.com/absaoss/k8s_crd.git
cd k8s_crd
go build cmd/coredns.go
./coredns -plugins | grep k8s_crd
```

For more details refer to [this CoreDNS doc](https://coredns.io/2017/07/25/compile-time-enabling-or-disabling-plugins/)
