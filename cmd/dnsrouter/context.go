package main

import (
	"net"

	"github.com/miekg/dns"
)

type Context struct {
}

// convenience function to return the first question
func (me *Context) Question(r *dns.Msg) dns.Question {
	return r.Question[0]
}

func (me *Context) WriteRR(w dns.ResponseWriter, r *dns.Msg, rr_string string) error {
	reply := new(dns.Msg)
	reply.SetReply(r)
	rr, err := dns.NewRR(rr_string)
	if err != nil {
		log.Warnf("NewRR %s: %s", rr_string, err)
		return err
	}
	reply.Answer = append(reply.Answer, rr)
	w.WriteMsg(reply)
	return nil
}

// return an A record reply message
func (me *Context) ReplyA(r *dns.Msg, name string, ip_string string, ttl uint32) *dns.Msg {
	ip := net.ParseIP(ip_string)
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
