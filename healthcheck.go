package dns_router
import (
   "fmt"
   "time"
)
const (
   StateUnknown = iota 
   StateUp
   StateDown
)
var CheckStates = map[CheckState]string {
   StateUnknown: "Unknown",
   StateUp: "Up",
   StateDown: "Down",
}
type CheckState int
func (s CheckState) String() string {
   return CheckStates[s]
}


type HealthChecker interface {
   Check([]CheckState) error
}


func HealthCheck(hc HealthChecker, backend Handler) {
   pstate := make([]CheckState, backend.ServerCount())
   state := make([]CheckState, backend.ServerCount())
   interval := time.Duration(5) * time.Second
   c := time.Tick(interval)

   fmt.Printf("healthcheck started\n")
   for {
      <- c
      err := hc.Check(state)
      if err != nil {
         fmt.Printf("Healthcheck backend error: %s", err)
         continue
      }

      for idx, alive := range state {
         if state[idx] != pstate[idx] {
            fmt.Printf("backend #%d %s changed state from %v to %v\n", idx, backend.ServerName(idx), pstate[idx], state[idx])
            backend.SetAlive(idx, alive)
         }
      }

      // save previous state
      for idx, alive := range state {
         pstate[idx]=alive
      }
   }
}

