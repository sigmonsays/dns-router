package dns_router

import (
	"fmt"
	"github.com/miekg/dns"
)

type Handler interface {
	SetAlive(idx int, alive CheckState) error
	ServerCount() int
	ServerName(idx int) string
}

type RequestHandler struct {

	// dns servers to use incase backend fails
	DefaultServers []string

	// healthcheck type
	HealthCheck HealthChecker

	// unique handler number
	Number int

	// where we send requests
	Servers []string
}

func NewRequestHandler(num_servers int, DefaultServers []string) *RequestHandler {
	h := &RequestHandler{
		DefaultServers: DefaultServers,
	}
	return h
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

func (h *RequestHandler) ServeDefaultDNS(w dns.ResponseWriter, r *dns.Msg) {
	addr := w.RemoteAddr()
	query := h.QueryDescription(r)
	c := new(dns.Client)
	for _, s := range h.DefaultServers {
		reply, _, err := c.Exchange(r, s)
		if err != nil {
			fmt.Printf("- [default] %s %d [%s] %s - %s\n", addr, r.Id, s, query, err)
			continue
		}
		fmt.Printf("+ [default] %s %d [%s] %s - %d\n", addr, r.Id, s, query, len(reply.Answer))
		w.WriteMsg(reply)
		break
	}
}
