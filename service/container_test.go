package service

import (
	"context"
	"testing"

	"github.com/coredns/coredns/plugin"
	"github.com/golang/mock/gomock"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockHandler(ctrl)
	w := NewMockResponseWriter(ctrl)
	tests := []struct {
		name     string
		handlers []struct {
			handler plugin.Handler
			f       *gomock.Call
			w       *gomock.Call
		}
	}{
		{
			name: "oneHandlerSucceed",
			handlers: []struct {
				handler plugin.Handler
				f       *gomock.Call
				w       *gomock.Call
			}{
				{
					handler: m,
					f:       m.EXPECT().ServeDNS(gomock.Any(), gomock.Any(), gomock.Any()).Return(dns.RcodeSuccess, nil),
					w:       w.EXPECT().WriteMsg(gomock.Any()).Return(nil),
				},
			},
		},
	}

	for _, test := range tests {
		c := NewCommonContainer()
		for _, h := range test.handlers {
			err := c.Register(h.handler)
			assert.NoError(t, err)
		}
		err := c.Execute(context.TODO(), w, &dns.Msg{})
		assert.NoError(t, err)
	}
}
