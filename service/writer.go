package service

import (
	"net"

	"github.com/miekg/dns"
)

// containerResponseWriter writer allows access to Message object (written by w.WriteMsg()) which is hidden by dns.ResponseWriter.
// The containerResponseWriter wraps any ResponseWriter and add getMsg function which provides Written message, so the
// dns.Msg is accessible within the container pipeline.
type containerResponseWriter struct {
	w   dns.ResponseWriter
	msg *dns.Msg
}

func newContainerWriter(w dns.ResponseWriter) *containerResponseWriter {
	return &containerResponseWriter{
		w: w,
	}
}

// LocalAddr returns the net.Addr of the server
func (c *containerResponseWriter) LocalAddr() net.Addr {
	return c.w.LocalAddr()
}

// RemoteAddr returns the net.Addr of the client that sent the current request.
func (c *containerResponseWriter) RemoteAddr() net.Addr {
	return c.w.RemoteAddr()
}

// WriteMsg saves a message that you can pick using getMessage().
// Function override original functionality and doesn't write anything to the response!
// Use this function as you use it with ResponseWriter. The container takes care about logic.
func (c *containerResponseWriter) WriteMsg(msg *dns.Msg) error {
	c.msg = msg
	return nil
}

// WriteContainerResult is equal to ResponseWriter.WriteMsg()
func (c *containerResponseWriter) WriteContainerResult() error {
	return c.w.WriteMsg(c.msg)
}

// Write writes a raw buffer back to the client.
func (c *containerResponseWriter) Write(bytes []byte) (int, error) {
	return c.w.Write(bytes)
}

// Close closes the connection.
func (c *containerResponseWriter) Close() error {
	return c.w.Close()
}

// TsigStatus returns the status of the Tsig.
func (c *containerResponseWriter) TsigStatus() error {
	return c.w.TsigStatus()
}

// TsigTimersOnly sets the tsig timers only boolean.
func (c *containerResponseWriter) TsigTimersOnly(b bool) {
	c.w.TsigTimersOnly(b)
}

// Hijack lets the caller take over the connection.
// After a call to Hijack(), the DNS package will not do anything with the connection.
func (c *containerResponseWriter) Hijack() {
	c.w.Hijack()
}

func (c *containerResponseWriter) getMsg() *dns.Msg {
	return c.msg
}
