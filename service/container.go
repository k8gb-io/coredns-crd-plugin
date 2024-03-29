package service

/*
Copyright 2022 The k8gb Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	"context"
	"fmt"

	clog "github.com/coredns/coredns/plugin/pkg/log"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

const thisPlugin = "container"

var log = clog.NewWithPlugin(thisPlugin)

type Container struct {
	services []plugin.Handler
}

func NewCommonContainer() *Container {
	return &Container{
		services: []plugin.Handler{},
	}
}

func (c *Container) Register(handler plugin.Handler) error {
	if handler == nil {
		return fmt.Errorf("nil plugin")
	}
	c.services = append(c.services, handler)
	return nil
}

func (c *Container) Execute(ctx context.Context, w dns.ResponseWriter, msg *dns.Msg) (err error) {
	var rcode int
	wr := newContainerWriter(w, msg)
	for _, svc := range c.services {
		rcode, err = svc.ServeDNS(ctx, wr, msg)
		if err != nil {
			return fmt.Errorf("%s: %w", svc.Name(), err)
		}
		if rcode != dns.RcodeSuccess {
			log.Errorf("skipping plugin %s", svc.Name())
			return wr.WriteContainerResult()
		}
		msg = wr.message()
	}
	return wr.WriteContainerResult()
}
