module github.com/AbsaOSS/k8s_crd

go 1.15

replace go.etcd.io/etcd => go.etcd.io/etcd v0.5.0-alpha.5.0.20200425165423-262c93980547

require (
	github.com/Azure/azure-sdk-for-go v46.0.0+incompatible // indirect
	github.com/coredns/caddy v1.1.0
	github.com/coredns/coredns v1.8.1
	github.com/miekg/dns v1.1.35
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	sigs.k8s.io/external-dns v0.7.6
)
