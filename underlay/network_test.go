package underlay

import (
  "testing"
  "math"
  "math/rand"
  "time"
)

func TestRandomUniformNetworkConnected(t *testing.T) {
  rand.Seed(time.Now().UTC().UnixNano())

  network := NewRandomUniformNetwork(10000, 70000, 2, 10)
  assertEqual(t, connected(network), true)
}

func TestRandomUniformNetworkDegreeStdReport(t *testing.T) {
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

  if rep > 0.5 {
    t.Fatalf("Too big coefficient of variation.")
  }
}

func TestGenerateTransitDomainGraphCorrectEdgeWeights(t *testing.T) {
  network := generateTransitDomainGraph(100, Wtt, Wttd)
  for _, node := range network.Routers {
    for _, conn := range node.Connections() {
      if conn.Latency() > Wtt + Wttd || conn.Latency() < Wtt - Wttd {
        t.Fatalf("Wrong edge weights.")
      }
    }
  }
}

// func TestTransitDomainsAreConnected(t *testing.T) {
//   tdg := generateTransitDomainGraph(100, Wtt, Wttd)
//   backbone := generateTransitDomains(tdg, 10)
//
//   assertEqual(t, connected(backbone), true)
// }
//
// func TestTransitDomainsHaveCorrectNumberOfNodes(t *testing.T) {
//   for tt := 0; tt < 5; tt++ {
//     doms := 20
//     domN := 20
//
//     tdg := generateTransitDomainGraph(doms, Wtt, Wttd)
//     backbone := generateTransitDomains(tdg, domN)
//
//     if len(backbone.Routers) < (domN - Ntd) * doms {
//       t.Fatalf("Backbone has too few nodes.")
//     }
//     if len(backbone.Routers) > (domN + Ntd) * doms {
//       t.Fatalf("Backbone has too many nodes.")
//     }
//   }
// }
