package service

import (
	"context"
	"fmt"

	"github.com/coredns/coredns/plugin"
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
	wr := newContainerWriter(w)
	for _, svc := range c.services {
		rcode, err = svc.ServeDNS(ctx, wr, msg)
		if wr.MessageWasWritten() {
			msg = wr.getMsg()
		}
		if err != nil {
			return fmt.Errorf("%s: %w", svc.Name(), err)
		}
		if rcode != dns.RcodeSuccess {
			return fmt.Errorf("[service: %s]", svc.Name())
		}
	}
	return nil
}
