package main
import (
   "fmt"
   "github.com/miekg/dns"
)

type PatternHandler struct {
   Pattern string

   // unique handler number
   Number int

   // where we send requests
   Servers []string
}

func (h *PatternHandler) QueryDescription(r *dns.Msg) string {
   return fmt.Sprintf("%s", r.Question[0].Name)
}
func (h *PatternHandler) AnswerDescription(r *dns.Msg) string {
   return fmt.Sprintf("%s", r.Answer)
}

func (h *PatternHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
   addr := w.RemoteAddr()
   query := h.QueryDescription(r)
   c := new(dns.Client)
   for _, s := range h.Servers {
      reply, _, err := c.Exchange(r, s)
      if err != nil {
         fmt.Printf("- %s %d [%s %s] %s - %s\n",  addr, r.Id, h.Pattern, s, query, err)
         continue
      }
      fmt.Printf("+ %s %d [%s %s] %s - %d\n", addr, r.Id, h.Pattern, s, query, len(reply.Answer))
      w.WriteMsg(reply)
      break
   }
}
