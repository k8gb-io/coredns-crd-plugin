package main

import (
	_ "github.com/AbsaOSS/k8s_crd"
	_ "github.com/coredns/coredns/core/plugin"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"
)

var dropPlugins = map[string]bool{
	"kubernetes":   true,
	"k8s_external": true,
}

func init() {
	var directives []string
	var alreadyAdded bool

	for _, name := range dnsserver.Directives {

		if dropPlugins[name] {
			if !alreadyAdded {
				directives = append(directives, "k8s_crd")
				alreadyAdded = true
			}
			continue
		}
		directives = append(directives, name)
	}

	dnsserver.Directives = directives

}

func main() {
	coremain.Run()
}
