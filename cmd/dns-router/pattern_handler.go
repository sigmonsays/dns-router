package main

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/sigmonsays/dns-router"
)

type PatternHandler struct {
	*dns_router.RequestHandler

	Pattern string
}

func (h *PatternHandler) AnswerDescription(r *dns.Msg) string {
	return fmt.Sprintf("%s", r.Answer)
}
func (h *PatternHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {

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
			fmt.Printf("- %s %d [%s %s] %s - %s\n", addr, r.Id, h.Pattern, s, query, err)
			continue
		}
		fmt.Printf("+ %s %d [%s %s] %s - %d\n", addr, r.Id, h.Pattern, s, query, len(reply.Answer))
		w.WriteMsg(reply)
		break
	}
}

func (h *PatternHandler) ServerCount() int {
	return len(h.Servers)
}
