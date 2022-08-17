package wrr

import (
	"context"
	"github.com/AbsaOSS/k8s_crd/service"
	"github.com/golang/mock/gomock"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO: combina A, AAA, MX CNAme records
type fakeWriter struct {
	w dns.ResponseWriter
}

func newFakeWriter(ctrl *gomock.Controller, f func(w *service.MockResponseWriter)) fakeWriter {
	w := service.NewMockResponseWriter(ctrl)
	f(w)
	return fakeWriter{
		w: w,
	}
}

func TestWeightRoundRobin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		msg           *dns.Msg
		writer        fakeWriter
		expectedError bool
	}{
		{
			name:          "Serve empty answers",
			msg:           &dns.Msg{},
			writer:        newFakeWriter(ctrl, func(w *service.MockResponseWriter) {}),
			expectedError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wrr := NewWeightRoundRobin()
			code, err := wrr.ServeDNS(context.TODO(), test.writer.w, test.msg)
			assert.Equal(t, dns.RcodeSuccess, code)
			assert.Equal(t, test.expectedError, err != nil)
		})
	}
}
