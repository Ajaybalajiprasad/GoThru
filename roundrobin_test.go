package loadbalancer

import "testing"

func TestRoundRobinBasic(t *testing.T) {
    servers := []*Server{
        NewServer("s1:80"),
        NewServer("s2:80"),
        NewServer("s3:80"),
    }
    lb := NewRoundRobin(servers)

    expected := []string{"s1:80", "s2:80", "s3:80", "s1:80", "s2:80", "s3:80"}
    for i, want := range expected {
        got := lb.Next()
        if got == nil {
            t.Fatalf("request %d: got nil, want %s", i, want)
        }
        if got.Address != want {
            t.Errorf("request %d: got %s, want %s", i, got.Address, want)
        }
    }
}

func TestRoundRobinSkipsDeadServer(t *testing.T) {
    servers := []*Server{
        NewServer("s1:80"),
        NewServer("s2:80"),
        NewServer("s3:80"),
    }
    lb := NewRoundRobin(servers)
    servers[1].Alive = false

    for i := 0; i < 6; i++ {
        got := lb.Next()
        if got == nil {
            t.Fatal("got nil, expected a live server")
        }
        if got.Address == "s2:80" {
            t.Error("got dead server s2, should have been skipped")
        }
    }
}

func TestRoundRobinEmpty(t *testing.T) {
    lb := NewRoundRobin([]*Server{})
    if lb.Next() != nil {
        t.Error("expected nil for empty pool")
    }
}

var _ Balancer = (*RoundRobin)(nil)
