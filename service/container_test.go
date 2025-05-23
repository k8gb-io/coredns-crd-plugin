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
	"testing"

	"github.com/k8gb-io/coredns-crd-plugin/common/mocks"

	"github.com/coredns/coredns/plugin"
	dns "github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type fakePlugin struct {
	handler plugin.Handler
}

type fakeWriter struct {
	writer dns.ResponseWriter
}

func newFakePlugin(ctrl *gomock.Controller, f func(h *mocks.MockHandler)) fakePlugin {
	h := mocks.NewMockHandler(ctrl)
	f(h)
	return fakePlugin{
		handler: h,
	}
}

func newFakeWriter(ctrl *gomock.Controller, f func(w *mocks.MockResponseWriter)) fakeWriter {
	w := mocks.NewMockResponseWriter(ctrl)
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
			name: "One handler succeed",
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).Return(nil).Times(1)
			}),
			plugin: []fakePlugin{
				newFakePlugin(ctrl, func(h *mocks.MockHandler) {
					h.EXPECT().ServeDNS(gomock.Any(), gomock.Any(), gomock.Any()).Return(dns.RcodeSuccess, nil).Times(1)
				})},
		},
		{
			name: "Two handlers succeed",
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).Return(nil).Times(1)
			}),
			plugin: []fakePlugin{
				newFakePlugin(ctrl, func(h *mocks.MockHandler) {
					h.EXPECT().ServeDNS(gomock.Any(), gomock.Any(), gomock.Any()).Return(dns.RcodeSuccess, nil).Times(1)
				}),
				newFakePlugin(ctrl, func(h *mocks.MockHandler) {
					h.EXPECT().ServeDNS(gomock.Any(), gomock.Any(), gomock.Any()).Return(dns.RcodeSuccess, nil).Times(1)
				})},
		},
		{
			name: "No handlers",
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).Return(nil).Times(1)
			}),
			plugin: []fakePlugin{},
		},
		{
			name:          "First handler error",
			expectedError: true,
			writer:        newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {}),
			plugin: []fakePlugin{
				newFakePlugin(ctrl, func(h *mocks.MockHandler) {
					h.EXPECT().Name().Return("fake").Times(1)
					h.EXPECT().ServeDNS(gomock.Any(), gomock.Any(), gomock.Any()).Return(dns.RcodeSuccess, fmt.Errorf("fake")).Times(1)
				}),
				newFakePlugin(ctrl, func(h *mocks.MockHandler) {})},
		},
		{
			name:          "Second handler error",
			expectedError: true,
			writer:        newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {}),
			plugin: []fakePlugin{
				newFakePlugin(ctrl, func(h *mocks.MockHandler) {
					h.EXPECT().ServeDNS(gomock.Any(), gomock.Any(), gomock.Any()).Return(dns.RcodeSuccess, nil).Times(1)
				}),
				newFakePlugin(ctrl, func(h *mocks.MockHandler) {
					h.EXPECT().Name().Return("fake").Times(1)
					h.EXPECT().ServeDNS(gomock.Any(), gomock.Any(), gomock.Any()).Return(dns.RcodeSuccess, fmt.Errorf("fake")).Times(1)
				})},
		},
		{
			name:          "First handler refused status",
			expectedError: false,
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).Return(nil).Times(1)
			}),
			plugin: []fakePlugin{
				newFakePlugin(ctrl, func(h *mocks.MockHandler) {
					h.EXPECT().Name().Return("fake").Times(1)
					h.EXPECT().ServeDNS(gomock.Any(), gomock.Any(), gomock.Any()).Return(dns.RcodeRefused, nil).Times(1)
				}),
				newFakePlugin(ctrl, func(h *mocks.MockHandler) {})},
		},
		{
			name:          "Second module fails and written message also fails",
			expectedError: true,
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).Return(fmt.Errorf("broken message")).Times(1)
			}),
			plugin: []fakePlugin{
				newFakePlugin(ctrl, func(h *mocks.MockHandler) {
					h.EXPECT().ServeDNS(gomock.Any(), gomock.Any(), gomock.Any()).Return(dns.RcodeSuccess, nil).Times(1)
				}),
				newFakePlugin(ctrl, func(h *mocks.MockHandler) {
					h.EXPECT().Name().Return("fake").Times(1)
					h.EXPECT().ServeDNS(gomock.Any(), gomock.Any(), gomock.Any()).Return(dns.RcodeServerFailure, nil).Times(1)
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

func TestContainerWriter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m1 := &dns.Msg{}
	m2 := &dns.Msg{}
	tests := []struct {
		name            string
		w               fakeWriter
		initialMesage   *dns.Msg
		writtenMessage  *dns.Msg
		expectedMessage *dns.Msg
		expectedError   bool
	}{
		{
			name:            "Written writtenMessage has not been changed",
			w:               newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {}),
			writtenMessage:  m1,
			initialMesage:   m1,
			expectedMessage: m1,
			expectedError:   false,
		},
		{
			name:            "Written writtenMessage has been changed",
			w:               newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {}),
			writtenMessage:  m1,
			initialMesage:   m2,
			expectedMessage: m2,
			expectedError:   false,
		},
		{
			name:            "Written writtenMessage is nil",
			w:               newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {}),
			initialMesage:   m1,
			writtenMessage:  nil,
			expectedMessage: m1,
			expectedError:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := newContainerWriter(test.w.writer, test.initialMesage)
			err := w.WriteMsg(test.writtenMessage)
			assert.Equal(t, test.expectedError, err != nil)
			assert.Equal(t, test.expectedMessage, w.message())
		})
	}
}

func TestRegister(t *testing.T) {
	c := NewCommonContainer()
	err := c.Register(nil)
	assert.Error(t, err)
}
