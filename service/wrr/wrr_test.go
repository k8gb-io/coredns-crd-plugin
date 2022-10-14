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

	"github.com/AbsaOSS/k8s_crd/common/netutils"

	"github.com/AbsaOSS/k8s_crd/common/k8sctrl"

	"github.com/AbsaOSS/k8s_crd/common/mocks"
	"github.com/coredns/coredns/plugin/test"
	"github.com/golang/mock/gomock"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

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

	rs1 := newRecordSet([]dns.RR{
		test.AAAA("alpha.cloud.example.com.		300	IN	AAAA		4001:a1:1014::89"),
		test.CNAME("beta.cloud.example.com.	300	IN	CNAME		beta.cloud.example.com."),
		test.AAAA("alpha.cloud.example.com.		300	IN	AAAA		4001:a1:1014::8a"),
		test.AAAA("alpha.cloud.example.com.		300	IN	AAAA		4001:a1:1014::8b"),
		test.AAAA("alpha.cloud.example.com.		300	IN	AAAA		4001:a1:1014::ff"),
		test.MX("alpha.cloud.example.com.	300	IN	MX		1	mxa-alpha.cloud.example.com."),
	}, []dns.RR{
		test.CNAME("beta.cloud.example.com.	300	IN	CNAME		beta.cloud.example.com."),
		test.MX("alpha.cloud.example.com.	300	IN	MX		1	mxa-alpha.cloud.example.com."),
		test.AAAA("alpha.cloud.example.com.		300	IN	AAAA		4001:a1:1014::89"),
		test.AAAA("alpha.cloud.example.com.		300	IN	AAAA		4001:a1:1014::8a"),
		test.AAAA("alpha.cloud.example.com.		300	IN	AAAA		4001:a1:1014::ff"),
		test.AAAA("alpha.cloud.example.com.		300	IN	AAAA		4001:a1:1014::8b"),
	}, map[string]string{"strategy": "roundrobin",
		"weight-eu-0-50": "4001:a1:1014::89",
		"weight-eu-1-50": "4001:a1:1014::8a",
		"weight-za-0-0":  "4001:a1:1014::8b",
		"weight-us-0-50": "4001:a1:1014::ff"})

	rs2 := newRecordSet([]dns.RR{
		test.AAAA("alpha.cloud.example.com.		300	IN	AAAA		4001:a1:1014::89"),
		test.AAAA("alpha.cloud.example.com.		300	IN	AAAA		4001:a1:1014::8b"),
		test.AAAA("alpha.cloud.example.com.		300	IN	AAAA		4001:a1:1014::ff"),
	}, []dns.RR{
		test.AAAA("alpha.cloud.example.com.		300	IN	AAAA		4001:a1:1014::89"),
		test.AAAA("alpha.cloud.example.com.		300	IN	AAAA		4001:a1:1014::ff"),
		test.AAAA("alpha.cloud.example.com.		300	IN	AAAA		4001:a1:1014::8b"),
	}, map[string]string{"strategy": "roundrobin",
		"weight-eu-0-50": "4001:a1:1014::89",
		"weight-za-0-0":  "4001:a1:1014::8b",
		"weight-us-0-50": "4001:a1:1014::ff"})

	rs3 := newRecordSet([]dns.RR{
		test.AAAA("alpha.cloud.example.com.		300	IN	AAAA		4001:a1:1014::89"),
		test.AAAA("alpha.cloud.example.com.		300	IN	AAAA		4001:a1:1014::8b"),
	}, []dns.RR{}, map[string]string{"strategy": "roundrobin",
		"weight-eu-0-50": "4001:a1:1014::89",
		"weight-za-0-0":  "4001:a1:1014::8b",
		"weight-us-0-50": "4001:a1:1014::ff"})

	rs4 := newRecordSet([]dns.RR{
		test.A("alpha.cloud.example.com.		300	IN	A		10.0.0.1"),
		test.CNAME("beta.cloud.example.com.	300	IN	CNAME		beta.cloud.example.com."),
		test.A("alpha.cloud.example.com.		300	IN	A		10.0.0.2"),
		test.A("alpha.cloud.example.com.		300	IN	A		10.10.0.1"),
		test.A("alpha.cloud.example.com.		300	IN	A		10.20.0.1"),
		test.MX("alpha.cloud.example.com.	300	IN	MX		1	mxa-alpha.cloud.example.com."),
	}, []dns.RR{
		test.CNAME("beta.cloud.example.com.	300	IN	CNAME		beta.cloud.example.com."),
		test.MX("alpha.cloud.example.com.	300	IN	MX		1	mxa-alpha.cloud.example.com."),
		test.A("alpha.cloud.example.com.		300	IN	A		10.0.0.1"),
		test.A("alpha.cloud.example.com.		300	IN	A		10.0.0.2"),
		test.A("alpha.cloud.example.com.		300	IN	A		10.20.0.1"),
		test.A("alpha.cloud.example.com.		300	IN	A		10.10.0.1"),
	}, map[string]string{"strategy": "roundrobin",
		"weight-eu-0-50": "10.0.0.1",
		"weight-eu-1-50": "10.0.0.2",
		"weight-za-0-0":  "10.10.0.1",
		"weight-us-0-50": "10.20.0.1"})

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
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).Return(nil).Times(1)
			}),
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
			expectedError: false,
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
			name: "Serve RR endpoint with valid weight 1-1 portions",
			msg: &dns.Msg{
				Answer: []dns.RR{
					test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"),
					test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.2"),
				},
			},
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).Return(nil).Times(1)
			}),
			expectedError: false,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Labels:  map[string]string{"strategy": "roundrobin", "weight-0-eu-1": "10.240.0.1", "weight-0-us-1": "10.240.0.2"},
					Targets: []string{"10.240.0.1"},
				}
			},
			rcode: dns.RcodeSuccess,
		},

		{
			name: "Serve RR endpoint with one 1 weight label",
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
					Labels:  map[string]string{"strategy": "roundrobin", "weight-eu-0-1": "10.240.0.1"},
					Targets: []string{"10.240.0.1"},
				}
			},
			rcode: dns.RcodeSuccess,
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
			name: "A records doesn't match with labels (count)",
			msg: &dns.Msg{
				Answer: []dns.RR{
					test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"),
					test.A("alpha.cloud.example.com.		300	IN	A			10.240.1.1"),
				},
			},
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).Return(nil).Times(0)
			}),
			expectedError: false,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Labels:  map[string]string{"strategy": "roundrobin", "weight-us-0-50": "10.240.0.1"},
					Targets: []string{"10.240.0.1"},
				}
			},
			rcode: dns.RcodeSuccess,
		},

		{
			name: "A records doesn't match with labels (values)",
			msg: &dns.Msg{
				Answer: []dns.RR{
					test.A("alpha.cloud.example.com.		300	IN	A			20.240.0.1"),
					test.A("alpha.cloud.example.com.		300	IN	A			20.240.1.1"),
				},
			},
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).Return(nil).Times(0)
			}),
			expectedError: false,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Labels:  map[string]string{"strategy": "roundrobin", "weight-eu-0-50": "10.240.0.1", "weight-us-0-50": "10.240.0.2"},
					Targets: []string{"10.240.0.1"},
				}
			},
			rcode: dns.RcodeSuccess,
		},

		{
			name: "A records doesn't match with labels (values2)",
			msg: &dns.Msg{
				Answer: []dns.RR{
					test.A("alpha.cloud.example.com.		300	IN	A			10.10.10.1"),
					test.A("alpha.cloud.example.com.		300	IN	A			10.10.10.2"),
				},
			},
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).Return(fmt.Errorf("broken writer")).Times(0)
			}),
			expectedError: false,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Labels:  map[string]string{"strategy": "roundrobin", "weight-eu-0-50": "10.240.0.1", "weight-us-0-50": "10.240.1.1"},
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
			expectedError: false,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Labels:  map[string]string{"strategy": "roundrobin", "weight-eu-0-100": "10.240.0.1"},
					Targets: []string{"10.240.0.1"},
				}
			},
			rcode: dns.RcodeServerFailure,
		},

		{
			name: "Asymetric group test",
			msg: &dns.Msg{
				Answer: []dns.RR{
					test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"),
					test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.2"),
					test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.3"),
					test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.4"),
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
					Labels: map[string]string{"strategy": "roundrobin",
						"weight-eu-0-50": "10.240.0.1",
						"weight-eu-1-50": "10.240.0.2",
						"weight-eu-2-50": "10.240.0.3",
						"weight-eu-3-50": "10.240.0.4",
						"weight-us-0-50": "10.240.1.1"},
					Targets: []string{"10.240.0.1"},
				}
			},
			rcode: dns.RcodeSuccess,
		},

		{
			name: "Handle AAAA records mixed with CNAME and MX records",
			msg: &dns.Msg{
				Answer: rs1.answer,
			},
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).DoAndReturn(rs1.checkExpected).Times(1)
			}),
			expectedError: false,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Labels:  rs1.labels,
					Targets: []string{"4001:a1:1014::89", "4001:a1:1014::8a", "4001:a1:1014::ff"},
				}
			},
			rcode: dns.RcodeSuccess,
		},

		{
			name: "Handle A records mixed with CNAME and MX records",
			msg: &dns.Msg{
				Answer: rs4.answer,
			},
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).DoAndReturn(rs4.checkExpected).Times(1)
			}),
			expectedError: false,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Labels:  rs4.labels,
					Targets: []string{"4001:a1:1014::89", "4001:a1:1014::8a", "4001:a1:1014::ff"},
				}
			},
			rcode: dns.RcodeSuccess,
		},

		{
			name: "Handle AAAA records without CNAME and MX records",
			msg: &dns.Msg{
				Answer: rs2.answer,
			},
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
				w.EXPECT().WriteMsg(gomock.Any()).DoAndReturn(rs2.checkExpected).Times(1)
			}),
			expectedError: false,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Labels:  rs2.labels,
					Targets: []string{"4001:a1:1014::89", "4001:a1:1014::8a", "4001:a1:1014::ff"},
				}
			},
			rcode: dns.RcodeSuccess,
		},

		{
			name: "Handle AAAA records with incomplete Answer section",
			msg: &dns.Msg{
				Answer: rs3.answer,
			},
			writer: newFakeWriter(ctrl, func(w *mocks.MockResponseWriter) {
			}),
			expectedError: false,
			lookup: func(indexKey string, clientIP net.IP) (result k8sctrl.LocalDNSEndpoint) {
				return k8sctrl.LocalDNSEndpoint{
					DNSName: host,
					Labels:  rs3.labels,
					Targets: []string{"4001:a1:1014::89", "4001:a1:1014::8a", "4001:a1:1014::ff"},
				}
			},
			rcode: dns.RcodeSuccess,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wrr := NewWeightRoundRobin()
			test.msg.Question = append(test.msg.Question, dns.Question{Qtype: dns.TypeA})
			k8sctrl.Resources.DNSEndpoint.Lookup = test.lookup
			code, err := wrr.ServeDNS(context.TODO(), test.writer.w, test.msg)
			assert.Equal(t, test.rcode, code)
			assert.Equal(t, test.expectedError, err != nil)
		})

	}
}

type recordset struct {
	answer   []dns.RR
	expected []dns.RR
	labels   map[string]string
}

func newRecordSet(answer, expected []dns.RR, labels map[string]string) *recordset {
	rs := new(recordset)
	rs.expected = expected
	rs.answer = answer
	rs.labels = labels
	return rs
}

func (rs *recordset) checkExpected(msg *dns.Msg) error {
	if len(msg.Answer) != len(rs.expected) {
		return fmt.Errorf("expecting answers")
	}
	oip, o, onoip := netutils.ParseAnswerSection(msg.Answer)
	eip, e, enoip := netutils.ParseAnswerSection(rs.expected)

	if len(oip) != len(eip) {
		return fmt.Errorf("%v %v", o, e)
	}
	if len(onoip) != len(enoip) {
		return fmt.Errorf("%v %v", enoip, onoip)
	}
	for k := range oip {
		if _, ok := eip[k]; !ok {
			return fmt.Errorf("%s not found", k)
		}
	}
	return nil
}
