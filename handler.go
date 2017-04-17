package dns_router

import (
	"bytes"
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/miekg/dns"
)

type Handler interface {
	SetAlive(idx int, alive CheckState) error
	ServerCount() int
	ServerName(idx int) string
}

func NewRequestHandler(num_servers int, DefaultServers []string) *RequestHandler {
	h := &RequestHandler{
		DefaultServers: DefaultServers,
	}
	return h
}

type RequestHandler struct {

	// the top level directory of the root dir. usually the directory of the config file
	RootDir string

	// dns servers to use incase backend fails
	DefaultServers []string

	// healthcheck type
	HealthCheck HealthChecker

	// unique handler number
	Number int

	// where we send requests
	Servers []string

	RedisPool *redis.Pool
}

func (h *RequestHandler) ServerName(idx int) string {
	return h.Servers[idx]
}

func (h *RequestHandler) SetAlive(idx int, alive CheckState) error {
	return h.HealthCheck.SetAlive(idx, alive)
}

func (h *RequestHandler) BackendAlive() bool {
	return h.HealthCheck.BackendAlive()
}

func (h *RequestHandler) QueryDescription(r *dns.Msg) string {
	return fmt.Sprintf("%s", r.Question[0].Name)
}

func (h *RequestHandler) ServeDefaultDNS(rlog *bytes.Buffer, w dns.ResponseWriter, r *dns.Msg) {
	fmt.Fprintf(rlog, "default=true ")
	// addr := w.RemoteAddr()
	// query := h.QueryDescription(r)
	c := new(dns.Client)
	for _, s := range h.DefaultServers {
		reply, _, err := c.Exchange(r, s)
		if err != nil {
			// fmt.Printf("- [default] %s %d [%s] %s - %s\n", addr, r.Id, s, query, err)
			continue
		}
		// fmt.Printf("+ [default] %s %d [%s] %s - %d\n", addr, r.Id, s, query, len(reply.Answer))
		w.WriteMsg(reply)
		break
	}
}
