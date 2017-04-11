package main

import (
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/miekg/dns"
	"github.com/sigmonsays/dns-router"
	"os"
)

var spewconf = spew.ConfigState{
	Indent:         "  ",
	DisableMethods: true,
	MaxDepth:       5,
}

func main() {
	conf := dns_router.Default()

	var configfile string
	flag.StringVar(&configfile, "config", "/etc/dns-router/config.yaml", "configuration file")
	flag.Parse()

	err := conf.LoadYaml(configfile)
	if err != nil {
		fmt.Printf("LoadYaml %s: %s", configfile, err)
		os.Exit(1)
	}

	conf.PrintConfig()

	mux := dns.NewServeMux()

	for n, b := range conf.Backends {
		num := n + 1

		request_handler := dns_router.NewRequestHandler(len(b.Servers), conf.Default.Servers)
		request_handler.Number = num
		request_handler.Servers = b.Servers

		var healthcheck dns_router.HealthChecker
		if b.HealthCheck {
			healthcheck = dns_router.NewPingHealthCheck(b.Servers)
		} else {
			healthcheck = dns_router.NewNullHealthCheck()
		}
		request_handler.HealthCheck = healthcheck

		t := &PatternHandler{
			RequestHandler: request_handler,
			Pattern:        b.Pattern,
		}

		if b.HealthCheck {
			go dns_router.HealthCheck(conf.HealthCheck, healthcheck, t)
		}

		mux.Handle(b.Pattern, t)
		fmt.Printf("pattern=%s servers=%s\n", b.Pattern, b.Servers)
	}

	healthcheck := dns_router.NewNullHealthCheck()
	request_handler := dns_router.NewRequestHandler(len(conf.Default.Servers), conf.Default.Servers)
	request_handler.Number = 0
	request_handler.Servers = conf.Default.Servers
	request_handler.HealthCheck = healthcheck
	t := &PatternHandler{
		RequestHandler: request_handler,
		Pattern:        ".",
	}
	mux.Handle(".", t)

	srv := &dns.Server{
		Addr:    conf.BindAddr,
		Net:     "udp",
		Handler: mux,
	}

	fmt.Printf("\nready\n")
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Printf("Listen %s\n", err)
	}

}
