package main

import (
	_ "github.com/AbsaOSS/k8s_crd"
	"github.com/AbsaOSS/k8s_crd/common/directives"
	_ "github.com/coredns/coredns/core/plugin"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"
)

func init() {
	p := directives.NewDirectivesManager(dnsserver.Directives)
	p.InsertBefore("k8s_crd", "kubernetes")
	p.Remove("kubernetes")
	p.Remove("k8s_external")
	dnsserver.Directives = p.Get()
}

func main() {
	coremain.Run()
}
