package dns_router
import (
   "fmt"
   "net"
   "sync"
   "os/exec"
)

type PingHealthCheck struct {
   mux sync.RWMutex
   Servers []string
   State []CheckState
}

func NewPingHealthCheck(servers []string) *PingHealthCheck {
   hc := &PingHealthCheck{
      Servers: servers,
      State: make([]CheckState, len(servers)),
   }
   for i := 0; i<len(servers); i++ {
      hc.State[i]=StateUnknown
   }
   return hc
}

func (hc *PingHealthCheck) Check(state []CheckState) error {
   for idx, s := range hc.Servers {
      host, _, err := net.SplitHostPort(s)
      if err != nil {
         fmt.Printf("healthcheck error %s: %s\n", s, err)
         continue
      }

      cmd := exec.Command("ping", "-n", "-c1", "-W", "1", host)
      err = cmd.Run()
      if err != nil {
         state[idx]=StateDown
      } else {
         state[idx]=StateUp
      }
   }
   return nil
}

func (hc *PingHealthCheck) BackendAlive() bool {

   hc.mux.RLock()
   defer hc.mux.RUnlock()

   // if any backend is alive we say the group is alive...
   for _, alive := range hc.State {
      if alive == StateUp {
         return true
      }
   }
   return false
}

func (hc *PingHealthCheck) SetAlive(idx int, alive CheckState) error {
   hc.mux.Lock()
   defer hc.mux.Unlock()
   hc.State[idx]=alive
   return nil
}
