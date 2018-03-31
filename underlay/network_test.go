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

func checkNetworkEdges(t *testing.T, network *Network, minL, maxL int) {
  for _, node := range network.Routers {
    for _, conn := range node.Connections() {
      if conn.Latency() > maxL || conn.Latency() < minL {
        t.Fatalf("Wrong edge weights.")
      }
    }
  }
}

func TestGenerateTransitDomainGraphCorrectEdgeWeights(t *testing.T) {
  network := generateTransitDomainGraph(100, Wtt, Wttd)
  checkNetworkEdges(t, network, Wtt - Wttd, Wtt + Wttd)
}

func TestNewRandomNetworkCorrectEdgeWeightsAndNodeNumber(t *testing.T) {
  for tt := 0; tt < 10; tt++ {
    N      := 10
    nodesD := 5
    edgeF  := 2
    minL   := 2
    maxL   := 10

    network := newRandomNetwork(N, 10, nodesD, edgeF, minL, maxL)
    checkNetworkEdges(t, network, minL, maxL)
    assertEqual(t, connected(network), true)

    N = 100
    nodesD = 10
    network = newRandomNetwork(N, 10, nodesD, edgeF, minL, maxL)
    checkNetworkEdges(t, network, minL, maxL)
    assertEqual(t, connected(network), true)
  }
}

func TestTransitDomainsAreConnected(t *testing.T) {
  tdg := generateTransitDomainGraph(2, Wtt, Wttd)
  backbone := generateTransitDomains(tdg, 5)
  assertEqual(t, connected(backbone), true)

  tdg = generateTransitDomainGraph(10, Wtt, Wttd)
  backbone = generateTransitDomains(tdg, 100)
  assertEqual(t, connected(backbone), true)

  tdg = generateTransitDomainGraph(100, Wtt, Wttd)
  backbone = generateTransitDomains(tdg, 10)
  assertEqual(t, connected(backbone), true)
}

func TestTransitDomainsHaveCorrectNumberOfNodes(t *testing.T) {
  for tt := 0; tt < 5; tt++ {
    doms := 15
    domN := 25

    tdg := generateTransitDomainGraph(doms, Wtt, Wttd)
    backbone := generateTransitDomains(tdg, domN)

    if len(backbone.Routers) < (domN - Ntd) * doms {
      t.Fatalf("Backbone has too few nodes.")
    }
    if len(backbone.Routers) > (domN + Ntd) * doms {
      t.Fatalf("Backbone has too many nodes.")
    }
  }
}

func TestAddedStubsAreConnected(t *testing.T) {
  doms  := 15
  stubs := 20
  stubN := 10

  backbone := generateTransitDomainGraph(doms, Wtt, Wttd)
  _, stubed := addStubs(backbone, stubs, stubN)
  assertEqual(t, connected(stubed), true)

  stubs = 1000
  backbone = generateTransitDomainGraph(doms, Wtt, Wttd)
  _, stubed = addStubs(backbone, stubs, stubN)
  assertEqual(t, connected(stubed), true)
}

func TestAddedStubsHaveCorrentNumberOfNodes(t *testing.T) {
  for tt := 0; tt < 5; tt++ {
    doms  := 15
    stubs := 20
    stubN := 10

    backbone := generateTransitDomainGraph(doms, Wtt, Wttd)
    _, stubed := addStubs(backbone, stubs, stubN)

    if len(stubed.Routers) < (stubN - Nsd) * stubs + doms {
      t.Fatalf("Stubed network has too few nodes.")
    }
    if len(stubed.Routers) > (stubN + Nsd) * stubs + doms {
      t.Fatalf("Stubed network too many nodes.")
    }
  }
}

func TestInternetworkConnected(t *testing.T) {
  network := NewInternetwork(15, 30, 200, 50)
  assertEqual(t, connected(network), true)

  network = NewInternetwork(2, 15, 100, 100)
  assertEqual(t, connected(network), true)
}
