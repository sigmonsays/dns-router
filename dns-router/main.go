package main
import (
   "os"
   "fmt"
   "flag"
   "github.com/miekg/dns"
   "github.com/sigmonsays/dns-router"
   "github.com/davecgh/go-spew/spew"
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

   for num, b := range conf.Backends {
      t := &PatternHandler{
         Pattern: b.Pattern,
         Number: num+1,
         Servers: b.Servers,
      }
      mux.Handle(b.Pattern, t)
      fmt.Printf("pattern=%s servers=%s\n", b.Pattern, b.Servers)
   }
   t := &PatternHandler{
      Pattern: ".",
      Number: 0,
      Servers: conf.Default.Servers,
   }
   mux.Handle(".", t)

   srv := &dns.Server{
      Addr: conf.BindAddr,
      Net: "udp",
      Handler: mux,
   }

   fmt.Printf("\nready\n")
   err = srv.ListenAndServe()
   if err != nil {
      fmt.Printf("Listen %s\n", err)
   }
   
}

