package loadbalancer

import "sync"

type RoundRobin struct {
  servers []*Server
  current int
  mu sync.Mutex
}

func NewRoundRobin(servers []*Server) *RoundRobin {
  return &RoundRobin{
    servers: servers,
    current: 0,
  }
}

func (rr* RoundRobin) Next() *Server{
  rr.mu.Lock()
  defer rr.mu.Unlock()

  n := len(rr.servers)
  if n == 0 {
    return nil
  }

  for i := 0; i<n; i++ {
    server := rr.servers[rr.current]
    rr.current = (rr.current + 1)%n
    if server.Alive {
      return server
    }
  }
  return nil
}

func (rr* RoundRobin) Add(server *Server){
  rr.mu.Lock()
  defer rr.mu.Unlock()
  rr.servers = append(rr.servers, server)
}

func (rr* RoundRobin) Servers() []*Server {
  rr.mu.Lock()
  defer rr.mu.Unlock()
  result := make([]*Server, len(rr.servers))
  copy(result, rr.servers)
  return result
}
