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

func (c *Container) Execute(ctx context.Context, w dns.ResponseWriter, msg *dns.Msg) error {
	for _, svc := range c.services {
		rcode, err := svc.ServeDNS(ctx, w, msg)
		if err != nil {
			return fmt.Errorf("%s: %w", svc.Name(), err)
		}
		if rcode != dns.RcodeSuccess {
			return fmt.Errorf("%s: returns unsuccesfull code", svc.Name())
		}
	}
	return nil
}