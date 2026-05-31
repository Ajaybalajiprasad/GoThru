package main

import (
  "fmt"
  "flag"
  "log"
  "net/http"
  "net/http/httputil"
  "net/url"
  "strings"
  "time"
  "github.com/Ajaybalajiprasad/GoThru"
)

func main(){
  port :=  flag.String("port", "8000", "port to run load balancer on")
  backends := flag.String("backends", "", "coma separated backedn url")
  flag.Parse()

  if *backends == "" {
    log.Fatal("please provide --backends flag")
  }

  var servers []*loadbalancer.Server
  for _, addr := range strings.Split(*backends, ",") {
    servers = append(servers, loadbalancer.NewServer(strings.TrimSpace(addr)))
  }

  lb := loadbalancer.NewRoundRobin(servers)

  go func() {
    for {
      for _, server := range lb.Servers() {
        resp, err := http.Get(server.Address + "/")
        if err != nil || resp.StatusCode != 200 {
          if server.Alive {
            log.Printf("server %s is DOWN", server.Address)
          }
          server.Alive = false
        }else {
          if !server.Alive {
            log.Printf("sever %s is UP!", server.Address)
          }
          server.Alive = true
        }
      }
      time.Sleep(10 * time.Second)
    }
  }()

  http.HandleFunc("/", func(w http.ResponseWriter, r * http.Request) {
    server := lb.Next()
    if server == nil {
      http.Error(w, "no servers available", 503)
      return
    }
    target, _ := url.Parse(server.Address)
    proxy := httputil.NewSingleHostReverseProxy(target)
    log.Printf("%s %s -> %s", r.Method, r.URL.Path, server.Address)
    proxy.ServeHTTP(w,r)
  })

  fmt.Printf("loadbalancer running on :%s\n", *port)
  fmt.Printf("Backends: %s\n", *backends)
  log.Fatal(http.ListenAndServe(":"+*port, nil))
}

