package k8scrd

import (
	"fmt"
	"strconv"

	"github.com/AbsaOSS/k8s_crd/service/gateway"
	"github.com/AbsaOSS/k8s_crd/service/wrr"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

type args struct {
	annotation     string
	apex           string
	filter         string
	kubecontroller string
	loadbalance    string
	negttl         uint32
	ttl            uint32
	zones          []string
}

const thisPlugin = "k8s_crd"
const weightRoundRobin = "weight"

var log = clog.NewWithPlugin(thisPlugin)

func init() {
	plugin.Register(thisPlugin, setup)
}

func setup(c *caddy.Controller) error {

	rawArgs, err := parse(c)
	if err != nil {
		return plugin.Error(thisPlugin, err)
	}

	k8sCRD, err := NewK8sCRD(configType(rawArgs.kubecontroller), rawArgs.filter)
	if err != nil {
		return plugin.Error(thisPlugin, err)
	}
	gwopts := gateway.NewGatewayOpts(rawArgs.annotation, rawArgs.apex, rawArgs.ttl, rawArgs.negttl, rawArgs.zones)
	_ = k8sCRD.container.Register(gateway.NewGateway(gwopts))
	if rawArgs.loadbalance == weightRoundRobin {
		_ = k8sCRD.container.Register(wrr.NewWeightRoundRobin())
	}
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		k8sCRD.Next = next
		return k8sCRD
	})

	return nil
}

func parseTTL(opt, arg string) (uint32, error) {
	t, err := strconv.Atoi(arg)
	if err != nil {
		return uint32(t), err
	}
	if t < 0 || t > 3600 {
		return uint32(t), fmt.Errorf("%s must be in range [0, 3600]: %d", opt, t)
	}
	return uint32(t), nil
}

func parse(c *caddy.Controller) (args, error) {
	a := args{}
	for c.Next() {
		a.zones = plugin.OriginsFromArgsOrServerBlock(c.RemainingArgs(), c.ServerBlockKeys)

		for c.NextBlock() {
			key := c.Val()
			args := c.RemainingArgs()
			if len(args) == 0 {
				return a, c.ArgErr()
			}
			switch key {
			case "filter":
				log.Infof("Filter: %+v", args)
				a.filter = args[0]
			case "annotation":
				log.Infof("annotation: %+v", args)
				a.annotation = args[0]
			case "ttl":
				ttl, err := parseTTL(c.Val(), args[0])
				if err != nil {
					a.ttl = ttl
				}
			case "negttl":
				log.Infof("negTTL: %+v", args[0])
				negttl, err := parseTTL(c.Val(), args[0])
				if err == nil {
					a.negttl = negttl
				}
			case "apex":
				a.apex = args[0]
			case "kubecontroller":
				log.Infof("kubecontroller: %+v", args)
				a.kubecontroller = args[0]
			case "loadbalance":
				log.Infof("loadbalance: %+v", args)
				a.loadbalance = args[0]
			default:
				return a, c.Errf("Unknown property '%s'", c.Val())
			}
		}
	}
	return a, nil
}
