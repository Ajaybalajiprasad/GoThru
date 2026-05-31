package loadbalancer

type Balancer interface {
  Next() *Server
  Add(server *Server)
  Servers() []*Server
}

func New(algorithm string, servers []*Server) Balancer {
  switch algorithm {
  case "weighted":
    return NewWeightedRoundRobin(servers)
  default:
    return NewRoundRobin(servers)
  }
} 
