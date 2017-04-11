package main

import "gopkg.in/natefinch/lumberjack.v2"

import (
	"github.com/miekg/dns"
)

func NewWrappedHandler(h dns.Handler, opts *WrappedOptions) *WrappedHandler {

	log := &lumberjack.Logger{
		Filename: opts.LogFile,
		// MaxSize:    500, // megabytes
		MaxBackups: 30,
		MaxAge:     120, // days
	}
	return &WrappedHandler{
		Log:     log,
		Handler: h,
	}
}

type WrappedHandler struct {
	Log     *lumberjack.Logger
	Handler dns.Handler
}
type WrappedOptions struct {
	LogFile string
}

func (me *WrappedHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {

	me.Handler.ServeDNS(w, r)

}
