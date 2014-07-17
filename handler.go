package dns_router
import (
   "fmt"
   "sync"
   "github.com/miekg/dns"
)

type Handler interface {
   SetAlive(idx int, alive CheckState) error
   ServerCount() int
   ServerName(idx int) string
}


type RequestHandler struct {
   mux sync.RWMutex

   // dns servers to use incase backend fails
   DefaultServers []string

   // healthcheck type
   HealthCheck HealthChecker

   // if the handler is alive according to health check
   alive []CheckState

   // unique handler number
   Number int

   // where we send requests
   Servers []string
}

func NewRequestHandler(num_servers int, DefaultServers []string) *RequestHandler {
   h := &RequestHandler{
      DefaultServers: DefaultServers,
      alive: make([]CheckState, num_servers+1),
   }
   return h
}

func (h *RequestHandler) ServerName(idx int) string {
   return h.Servers[idx]
}

func (h *RequestHandler) SetAlive(idx int, alive CheckState) error {
   h.mux.Lock()
   defer h.mux.Unlock()
   h.alive[idx]=alive
   return nil
}

func (h *RequestHandler) BackendAlive() bool {
   h.mux.RLock()
   defer h.mux.RUnlock()

   // if any backend is alive we say the group is alive...
   for _, alive := range h.alive {
      if alive == StateUp {
         return true
      }
   }
   return false
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
         fmt.Printf("- [default] %s %d [%s] %s - %s\n",  addr, r.Id, s, query, err)
         continue
      }
      fmt.Printf("+ [default] %s %d [%s] %s - %d\n", addr, r.Id, s, query, len(reply.Answer))
      w.WriteMsg(reply)
      break
   }
}
