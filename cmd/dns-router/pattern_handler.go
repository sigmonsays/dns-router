package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/miekg/dns"
	"github.com/sigmonsays/dns-router"
)

type PatternHandler struct {
	*dns_router.RequestHandler

	Override map[string]string
	IPAlias  map[string]string
	Log      io.Writer
	Pattern  string
}

func (h *PatternHandler) AnswerDescription(r *dns.Msg) string {
	return fmt.Sprintf("%s", r.Answer)
}

type OverrideResponse struct {

	// name found in the question section for this override
	Name string

	// found a valid override to return
	Found bool

	// a canned ip response
	IP string
}

func (h *PatternHandler) FindOverride(in *dns.Msg) *OverrideResponse {
	ret := &OverrideResponse{
		Found: false,
	}

	for _, q := range in.Question {
		if q.Qtype == dns.TypeA || q.Qtype == dns.TypeCNAME {
			result, ok := h.Override[q.Name]
			if ok {
				ret.Found = true
				ret.Name = q.Name
				ret.IP = result
				break
			}
		}
	}
	return ret
}

func (h *PatternHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {

	log := bytes.NewBuffer(nil) // request log

	override := h.FindOverride(r)
	if override.Found {

		fmt.Fprintf(log, "override=true ")

		reply := new(dns.Msg)
		reply.SetReply(r)
		rheader := dns.RR_Header{
			Name:   override.Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    300,
		}
		rr := new(dns.A)
		rr.Hdr = rheader
		rr.A = net.ParseIP(override.IP)
		reply.Answer = append(reply.Answer, rr)
		w.WriteMsg(reply)
		h.LogRoundTrip(log, w, r, reply)
		return
	}

	if h.BackendAlive() == false {
		h.ServeDefaultDNS(w, r)
		return
	}

	addr := w.RemoteAddr()
	query := h.QueryDescription(r)
	c := new(dns.Client)
	for _, s := range h.Servers {
		reply, _, err := c.Exchange(r, s)
		if err != nil {
			fmt.Printf("ERROR addr=%s id=%d [%s %s] query=%s: error %s\n", addr, r.Id, h.Pattern, s, query, err)
			continue
		}
		// fmt.Printf("+ %s %d [%s %s] %s - %d\n", addr, r.Id, h.Pattern, s, query, len(reply.Answer))
		w.WriteMsg(reply)
		h.LogRoundTrip(log, w, r, reply)
		break
	}
}

func (h *PatternHandler) ServerCount() int {
	return len(h.Servers)
}

func (h *PatternHandler) LogRoundTrip(log *bytes.Buffer, w dns.ResponseWriter, in *dns.Msg, out *dns.Msg) {
	laddr := w.LocalAddr()
	raddr := w.RemoteAddr()

	rhost, _, _ := net.SplitHostPort(raddr.String())
	rip := net.ParseIP(rhost)

	lhost, _, _ := net.SplitHostPort(laddr.String())
	lip := net.ParseIP(lhost)

	ralias, ralias_exists := h.IPAlias[rip.String()]
	lalias, lalias_exists := h.IPAlias[lip.String()]

	// output a single log record for a request/response

	//  0 src address
	//  1 dst address
	//  2 question
	//  3 answer

	buf := bytes.NewBuffer(nil)

	// build question variable
	for _, q := range in.Question {
		fmt.Fprintf(buf, "name=%s ", q.Name)
	}
	question := buf.String()

	// build answer variable
	buf.Reset()
	for _, rr := range out.Answer {
		fmt.Fprintf(buf, "%s ", FormatRR(rr))
	}
	answer := buf.String()

	buf.Reset()
	if ralias_exists == false {
		ralias = ""
	}
	fmt.Fprintf(buf, "%s/%s ", ralias, rip)
	if lalias_exists == false {
		lalias = ""
	}
	fmt.Fprintf(buf, "%s/%s ", lalias, lip)
	fmt.Fprintf(buf, "%s %s ", question, answer)

	// add additional log info
	buf.Write(log.Bytes())

	fmt.Printf("RoundTrip %s\n", buf.String())

	if h.Log != nil {
		fmt.Fprintf(h.Log, "%s %s\n", time.Now().Format(time.RFC3339Nano), buf.String())
	}

}

func FormatRR(rr dns.RR) string {

	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "%T ", rr)
	switch rr.(type) {

	case *dns.PTR:
		v, _ := rr.(*dns.PTR)
		fmt.Fprintf(buf, "ptr=%s", v.Ptr)

	case *dns.A:
		v, _ := rr.(*dns.A)
		fmt.Fprintf(buf, "ip=%s", v.A)

	case *dns.AAAA:
		v, _ := rr.(*dns.AAAA)
		fmt.Fprintf(buf, "ip=%s", v.AAAA)

	case *dns.CNAME:
		v, _ := rr.(*dns.CNAME)
		fmt.Fprintf(buf, "target=%s", v.Target)

	default:
		fmt.Fprintf(buf, "%#v", rr)
	}
	return buf.String()

}
