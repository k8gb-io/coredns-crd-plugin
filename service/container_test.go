package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/coredns/coredns/plugin"
	"github.com/golang/mock/gomock"
	dns "github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

type fakePlugin struct {
	handler plugin.Handler
}

type fakeWriter struct {
	writer dns.ResponseWriter
}

func newFakeHandler(ctrl *gomock.Controller, f func(h *MockHandler)) fakePlugin {
	h := NewMockHandler(ctrl)
	f(h)
	return fakePlugin{
		handler: h,
	}
}

func newFakeWriter(ctrl *gomock.Controller, f func(w *MockResponseWriter)) fakeWriter {
	w := NewMockResponseWriter(ctrl)
	f(w)
	return fakeWriter{
		writer: w,
	}
}

func TestContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		expectedError bool
		writer        fakeWriter
		plugin        []fakePlugin
	}{
		{
			name: "OneHandlerSucceed",
			writer: newFakeWriter(ctrl, func(w *MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).Return(nil).Times(1)
			}),
			plugin: []fakePlugin{
				newFakeHandler(ctrl, func(h *MockHandler) {
					h.EXPECT().ServeDNS(gomock.Any(), gomock.Any(), gomock.Any()).Return(dns.RcodeSuccess, nil).Times(1)
				})},
		},
		{
			name: "TwoHandlersSucceed",
			writer: newFakeWriter(ctrl, func(w *MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).Return(nil).Times(1)
			}),
			plugin: []fakePlugin{
				newFakeHandler(ctrl, func(h *MockHandler) {
					h.EXPECT().ServeDNS(gomock.Any(), gomock.Any(), gomock.Any()).Return(dns.RcodeSuccess, nil).Times(1)
				}),
				newFakeHandler(ctrl, func(h *MockHandler) {
					h.EXPECT().ServeDNS(gomock.Any(), gomock.Any(), gomock.Any()).Return(dns.RcodeSuccess, nil).Times(1)
				})},
		},
		{
			name: "NoHandlers",
			writer: newFakeWriter(ctrl, func(w *MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).Return(nil).Times(1)
			}),
			plugin: []fakePlugin{},
		},
		{
			name:          "FirstHandlerError",
			expectedError: true,
			writer:        newFakeWriter(ctrl, func(w *MockResponseWriter) {}),
			plugin: []fakePlugin{
				newFakeHandler(ctrl, func(h *MockHandler) {
					h.EXPECT().Name().Return("fake").Times(1)
					h.EXPECT().ServeDNS(gomock.Any(), gomock.Any(), gomock.Any()).Return(dns.RcodeSuccess, fmt.Errorf("fake")).Times(1)
				}),
				newFakeHandler(ctrl, func(h *MockHandler) {})},
		},
		{
			name:          "SecondHandlerError",
			expectedError: true,
			writer:        newFakeWriter(ctrl, func(w *MockResponseWriter) {}),
			plugin: []fakePlugin{
				newFakeHandler(ctrl, func(h *MockHandler) {
					h.EXPECT().ServeDNS(gomock.Any(), gomock.Any(), gomock.Any()).Return(dns.RcodeSuccess, nil).Times(1)
				}),
				newFakeHandler(ctrl, func(h *MockHandler) {
					h.EXPECT().Name().Return("fake").Times(1)
					h.EXPECT().ServeDNS(gomock.Any(), gomock.Any(), gomock.Any()).Return(dns.RcodeSuccess, fmt.Errorf("fake")).Times(1)
				})},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := NewCommonContainer()
			for _, p := range test.plugin {
				err := c.Register(p.handler)
				assert.NoError(t, err)
			}
			err := c.Execute(context.TODO(), test.writer.writer, &dns.Msg{})
			assert.Equal(t, test.expectedError, err != nil)
		})
	}
}
