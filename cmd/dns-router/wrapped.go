package main

import (
	"github.com/miekg/dns"
)

func NewWrappedHandler(h dns.Handler, opts *WrappedOptions) *WrappedHandler {

	return &WrappedHandler{
		Handler: h,
	}
}

type WrappedHandler struct {
	Handler dns.Handler
}
type WrappedOptions struct {
}

func (me *WrappedHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {

	me.Handler.ServeDNS(w, r)

}
