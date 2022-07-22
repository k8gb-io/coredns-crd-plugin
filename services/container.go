package services

import (
	"context"
	"fmt"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

// ServiceContainer is keeps particular services in the container and executes them
// in order they were added
type ServiceContainer interface {
	Add(plugin.Handler) error
	// Execute run plugin.Handler.ServeDNS in order as it was added. If error occurs,
	// it exits immediatelly end return error.
	Execute(ctx context.Context, w dns.ResponseWriter, msg *dns.Msg) error
}

type Container struct {
	services []plugin.Handler
}

func NewContainer() *Container {
	return &Container{
		services: []plugin.Handler{},
	}
}

func (c *Container) Add(handler plugin.Handler) error {
	if handler == nil {
		return fmt.Errorf("nil plugin")
	}
	c.services = append(c.services,handler)
	return nil
}

func (c *Container) Execute(ctx context.Context, w dns.ResponseWriter, msg *dns.Msg) (err error) {
	for _, svc := range c.services {
		_, err = svc.ServeDNS(ctx,w,msg)
		if err != nil {
			return err
		}
	}
	return nil
}
