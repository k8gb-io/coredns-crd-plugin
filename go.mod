module github.com/AbsaOSS/k8s_crd

go 1.16

replace go.etcd.io/etcd => go.etcd.io/etcd v0.5.0-alpha.5.0.20200425165423-262c93980547

require (
	github.com/coredns/caddy v1.1.1
	github.com/coredns/coredns v1.8.6
	github.com/maxmind/mmdbwriter v0.0.0-20210819141656-efe6d8ec5816
	github.com/miekg/dns v1.1.43
	github.com/oschwald/maxminddb-golang v1.8.0
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	sigs.k8s.io/external-dns v0.10.0
)
