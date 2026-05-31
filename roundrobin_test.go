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

func TestWeightedDistribution(t *testing.T) {
    servers := []*Server{
        NewWeightedServer("s1:80", 3),
        NewWeightedServer("s2:80", 2),
        NewWeightedServer("s3:80", 1),
    }
    lb := NewWeightedRoundRobin(servers)

    counts := map[string]int{}
    for i := 0; i < 60; i++ {
        s := lb.Next()
        counts[s.Address]++
    }

    if counts["s1:80"] != 30 {
        t.Errorf("s1 got %d, want 30", counts["s1:80"])
    }
    if counts["s2:80"] != 20 {
        t.Errorf("s2 got %d, want 20", counts["s2:80"])
    }
    if counts["s3:80"] != 10 {
        t.Errorf("s3 got %d, want 10", counts["s3:80"])
    }
}

var _ Balancer = (*WeightedRoundRobin)(nil)
