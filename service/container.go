package service

import (
	"context"
	"fmt"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/miekg/dns"
)

type Container struct {
	services []plugin.Handler
}

func NewCommonContainer() *Container {
	return &Container{
		services: []plugin.Handler{},
	}
}

func (c *Container) Add(handler plugin.Handler) error {
	if handler == nil {
		return fmt.Errorf("nil plugin")
	}
	c.services = append(c.services, handler)
	return nil
}

func (c *Container) Execute(ctx context.Context, w dns.ResponseWriter, msg *dns.Msg) (err error) {
	var rcode int
	for _, svc := range c.services {
		rcode, err = svc.ServeDNS(ctx, w, msg)
		msg = getMsg(w)
		if err != nil {
			return fmt.Errorf("%s: %w", svc.Name(), err)
		}
		if rcode != dns.RcodeSuccess {
			return fmt.Errorf("%s: returns unsuccesfull code", svc.Name())
		}
	}
	return nil
}

// getMsg reads written msg from dns.ResponseWriter. The message can be further modified in the next iteration by
//the container
// TODO: hack, casting interface into original interface. Consider reimplement gateway into dns.ResponseWriter.
func getMsg(w dns.ResponseWriter) (msg *dns.Msg) {
	impl :=  w.(*dnstest.Recorder)
	return impl.Msg
}