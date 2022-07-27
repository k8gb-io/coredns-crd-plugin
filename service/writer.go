package service

import (
	"github.com/miekg/dns"
	"net"
)

// container writer allows access to Message object (written by w.WriteMsg()) which is hidden by dns.ResponseWriter
// the struct is private, no access outside of package
type containerWriter struct {
	w dns.ResponseWriter
	msg *dns.Msg
}

func newContainerWriter(w dns.ResponseWriter) *containerWriter {
	return &containerWriter{
		w: w,
	}
}

// LocalAddr returns the net.Addr of the server
func (c *containerWriter) LocalAddr() net.Addr {
	return c.w.LocalAddr()
}
// RemoteAddr returns the net.Addr of the client that sent the current request.
func (c *containerWriter) RemoteAddr() net.Addr {
	return c.w.RemoteAddr()
}

// WriteMsg writes a reply back to the client.
func (c *containerWriter) WriteMsg(msg *dns.Msg) error {
	c.msg = msg
	return c.w.WriteMsg(msg)
}

// Write writes a raw buffer back to the client.
func (c *containerWriter) Write(bytes []byte) (int, error) {
	return c.w.Write(bytes)
}

// Close closes the connection.
func (c *containerWriter) Close() error {
	return c.w.Close()
}

// TsigStatus returns the status of the Tsig.
func (c *containerWriter) TsigStatus() error {
	return c.w.TsigStatus()
}

// TsigTimersOnly sets the tsig timers only boolean.
func (c *containerWriter) TsigTimersOnly(b bool) {
	c.w.TsigTimersOnly(b)
}

// Hijack lets the caller take over the connection.
// After a call to Hijack(), the DNS package will not do anything with the connection.
func (c *containerWriter) Hijack() {
	c.w.Hijack()
}

func (c *containerWriter) getMsg() *dns.Msg {
	return c.msg
}