package wrr

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
	"net"
	"testing"

	"github.com/AbsaOSS/k8s_crd/common/k8sctrl"

	"github.com/AbsaOSS/k8s_crd/common/mocks"
	"github.com/coredns/coredns/plugin/test"
	"github.com/golang/mock/gomock"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

// TODO: combine A, AAA, MX CNAme records
// TODO: endpoint values!
type fakeWriter struct {
	w dns.ResponseWriter
}

func newFakeWriter(ctrl *gomock.Controller, f func(w *mocks.MockResponseWriter)) fakeWriter {
	w := mocks.NewMockResponseWriter(ctrl)
	f(w)
	return fakeWriter{
		w: w,
	}
}

func TestWeightRoundRobin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const host = "roundrobin.cloud.example.com"

	tests := []struct {
		name          string
		msg           *dns.Msg
		writer        fakeWriter
		expectedError bool
		rcode         int
		lookup        k8sctrl.LookupEndpoint
	}{

		{
			name:          "Serve empty answers",
			msg:           &dns.Msg{},
			writer:        newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {}),
			expectedError: false,
		},

		{
			name: "Serve RR endpoint without labels",
			msg: &dns.Msg{
				Answer: []dns.RR{
					test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"),
				},
			},
			writer:        newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {}),
			expectedError: false,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Targets: []string{"10.0.0.1", "10.0.0.2"},
				}
			},
			rcode: dns.RcodeSuccess,
		},

		{
			name: "Serve RR endpoint without weight labels",
			msg: &dns.Msg{
				Answer: []dns.RR{
					test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"),
				},
			},
			writer:        newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {}),
			expectedError: false,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Labels:  map[string]string{"strategy": "roundrobin"},
					Targets: []string{"10.240.0.1"},
				}
			},
			rcode: dns.RcodeSuccess,
		},

		{
			name: "Serve RR endpoint with invalid weight label",
			msg: &dns.Msg{
				Answer: []dns.RR{
					test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"),
				},
			},
			writer:        newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {}),
			expectedError: true,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Labels:  map[string]string{"strategy": "roundrobin", "weight-0-eu": "10.240.0.1"},
					Targets: []string{"10.240.0.1"},
				}
			},
			rcode: dns.RcodeServerFailure,
		},

		{
			name: "Serve RR endpoint with invalid weight label2",
			msg: &dns.Msg{
				Answer: []dns.RR{
					test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"),
				},
			},
			writer:        newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {}),
			expectedError: true,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Labels:  map[string]string{"strategy": "roundrobin", "weight-0-eu-50": "10.240.0.1"},
					Targets: []string{"10.240.0.1"},
				}
			},
			rcode: dns.RcodeServerFailure,
		},

		{
			name: "Serve RR endpoint with one 50% weight label",
			msg: &dns.Msg{
				Answer: []dns.RR{
					test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"),
				},
			},
			writer:        newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {}),
			expectedError: true,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Labels:  map[string]string{"strategy": "roundrobin", "weight-eu-0-50": "10.240.0.1"},
					Targets: []string{"10.240.0.1"},
				}
			},
			rcode: dns.RcodeServerFailure,
		},

		{
			name: "Serve RR endpoint with one 100% weight label",
			msg: &dns.Msg{
				Answer: []dns.RR{
					test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"),
				},
			},
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).Return(nil).Times(1)
			}),
			expectedError: false,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Labels:  map[string]string{"strategy": "roundrobin", "weight-eu-0-100": "10.240.0.1"},
					Targets: []string{"10.240.0.1"},
				}
			},
			rcode: dns.RcodeSuccess,
		},

		{
			name: "Serve two 50% endpoints",
			msg: &dns.Msg{
				Answer: []dns.RR{
					test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"),
					test.A("alpha.cloud.example.com.		300	IN	A			10.240.1.1"),
				},
			},
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).Return(nil).Times(1)
			}),
			expectedError: false,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Labels:  map[string]string{"strategy": "roundrobin", "weight-us-0-50": "10.240.0.1", "weight-eu-0-50": "10.240.1.1"},
					Targets: []string{"10.240.0.1"},
				}
			},
			rcode: dns.RcodeSuccess,
		},

		{
			name: "Serve RR endpoint with one 100% weight label but broken writer",
			msg: &dns.Msg{
				Answer: []dns.RR{
					test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"),
				},
			},
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).Return(fmt.Errorf("broken writer")).Times(1)
			}),
			expectedError: true,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Labels:  map[string]string{"strategy": "roundrobin", "weight-eu-0-100": "10.240.0.1"},
					Targets: []string{"10.240.0.1"},
				}
			},
			rcode: dns.RcodeSuccess,
		},

		{
			name: "Serve RR endpoint where address doesn't meet weight label value",
			msg: &dns.Msg{
				Answer: []dns.RR{
					test.A("alpha.cloud.example.com.		300	IN	A			10.10.10.1"),
				},
			},
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).Return(fmt.Errorf("broken writer")).Times(1)
			}),
			expectedError: true,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Labels:  map[string]string{"strategy": "roundrobin", "weight-eu-0-100": "10.240.0.1"},
					Targets: []string{"10.240.0.1"},
				}
			},
			rcode: dns.RcodeSuccess,
		},
	}

	for _, unit := range tests {
		t.Run(unit.name, func(t *testing.T) {
			wrr := NewWeightRoundRobin()
			k8sctrl.Resources.DNSEndpoint.Lookup = unit.lookup
			code, err := wrr.ServeDNS(context.TODO(), unit.writer.w, unit.msg)
			assert.Equal(t, unit.rcode, code)
			assert.Equal(t, unit.expectedError, err != nil)
		})

	}
}
