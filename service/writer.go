package service

import (
	"net"

	"github.com/miekg/dns"
)

// containerResponseWriter writer allows access to Message object (written by w.WriteMsg()) which is hidden by dns.ResponseWriter.
// The containerResponseWriter wraps any ResponseWriter and add getMsg function which provides Written message, so the
// dns.Msg is accessible within the container pipeline.
type containerResponseWriter struct {
	w          dns.ResponseWriter
	msg        *dns.Msg
	wasWritten bool
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

// WriteMsg writes a reply back to the client.
func (c *containerResponseWriter) WriteMsg(msg *dns.Msg) error {
	c.msg = msg
	err := c.w.WriteMsg(msg)
	c.wasWritten = err == nil
	return err
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

// MessageWasWritten decides whether message was successfully written by ResponseWriter
func (c *containerResponseWriter) MessageWasWritten() bool {
	return c.wasWritten
}
