package main

import (
	"net"

	"github.com/miekg/dns"
)

type Context struct {
}

func (me *Context) ReplyA(r *dns.Msg, name string, ip net.IP, ttl uint32) *dns.Msg {
	reply := new(dns.Msg)
	reply.SetReply(r)
	rheader := dns.RR_Header{
		Name:   name,
		Rrtype: dns.TypeA,
		Class:  dns.ClassINET,
		Ttl:    ttl,
	}
	rr := new(dns.A)
	rr.Hdr = rheader
	rr.A = ip
	reply.Answer = append(reply.Answer, rr)
	return reply
}
