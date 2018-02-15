package underlay

import (
  "testing"
  "math"
  "math/rand"
  "time"
)

func TestRadondomUniformNetworkConnected(t *testing.T) {
  rand.Seed(time.Now().UTC().UnixNano())

  network := NewRandomUniformNetwork(10000, 70000, 2, 10)
  assertEqual(t, connected(network), true)
}

func TestRadondomUniformNetworkDegreeStdReport(t *testing.T) {
  rand.Seed(time.Now().UTC().UnixNano())
  network := NewRandomUniformNetwork(10000, 70000, 2, 10)

  mean := 0.0
  for i := range network.Routers {
    mean += float64(len(network.Routers[i].Connections()))
  }
  mean /= float64(len(network.Routers))

  variance := 0.0
  for i := range network.Routers {
    variance += math.Pow(float64(len(network.Routers[i].Connections())) - mean, 2)
  }
  variance /= float64(len(network.Routers))
  std := math.Sqrt(variance)
  rep := std/mean

  if rep > 0.3 {
    t.Fatalf("Too big coefficient of variation.")
  }
}
