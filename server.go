package loadbalancer

type Server struct {
  Address string
  Weight int
  Alive bool
}

func NewServer(address string) *Server {
  return &Server{
    Address: address,
    Weight: 1,
    Alive: true,
  }
}

func NewWeightedServer(address string, weight int) *Server{
  return &Server{
    Address: address,
    Weight: weight,
    Alive: true,
  }
}
