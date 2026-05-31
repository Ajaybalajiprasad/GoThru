package loadbalancer

import "sync"

type WeightedRoundRobin struct {
	servers []*Server
	sequence []*Server
	current int
	mu sync.Mutex
}

func NewWeightedRoundRobin(servers []*Server) *WeightedRoundRobin {
  wrr := &WeightedRoundRobin{
    servers: servers,
  }
  wrr.buildSequence()
  return wrr
}

func (wrr *WeightedRoundRobin) buildSequence() {
  wrr.sequence = []*Server{}
  for _, server := range wrr.servers {
    for i := 0; i<server.Weight; i++ {
      wrr.sequence = append(wrr.sequence, server)
    }
  }
  wrr.current = 0;
}

func (wrr *WeightedRoundRobin) Next() *Server {
  wrr.mu.Lock()
  defer wrr.mu.Unlock()

  n := len(wrr.sequence)
  if n == 0{
    return nil
  }

  for i := 0; i < n; i++ {
    server := wrr.sequence[wrr.current]
    wrr.current = (wrr.current + 1) % n
    if server.Alive {
      return server
    }
  }
  return nil
}

func (wrr *WeightedRoundRobin) Add(server *Server) {
  wrr.mu.Lock()
  defer wrr.mu.Unlock()
  wrr.servers = append(wrr.servers, server)
  wrr.buildSequence()
}

func (wrr *WeightedRoundRobin) Servers() []*Server {
  wrr.mu.Lock()
  wrr.mu.Unlock()
  result := make([]*Server, len(wrr.servers))
  copy(result, wrr.servers)
  return result
} 
