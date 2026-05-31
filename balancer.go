package loadbalancer

type Balancer interface {
  Next() *Server
  Add(server *Server)
  Servers() []*Server
}
