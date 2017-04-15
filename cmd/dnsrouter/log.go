package main

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("dnsrouter", func(newlog gologging.Logger) { log = newlog })
}
