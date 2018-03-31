package underlay

import (
  "math"
  "math/rand"
)

type Network struct {
  Routers []Router
}

func (n *Network) RandomRouter() Router {
  return n.Routers[rand.Intn(len(n.Routers))]
}

const Wtt  int = 100
const Wttd int = 2

const Ntd     int = 5
const minNt   int = 5
const edgeNtf int = 2

const minLatency int = 2
const maxLatency int = 10

/* Generates a 2-level transit-stub topology following the paper:
 Zegura E., Calvert K. and Bhattacharjee S. How to model an internetwork. In INFOCOM’96 (1996)

  T  = number of tranasit domains
  Nt = avg. nodes per transit domain
  S  = number of stub domains
  Ns = avg. nodes per stub domain

  The transit domains represent the backbone, whereas a stub is attached to a
  node in the backbone.

  The graph is generated as follows:
  1. Generate transit domain graph
  2. Generate graph from tranasit domain
  3. Add stubs
  4. Generate multi-homed stubs
*/
func NewInternetwork(T, Nt, S, Ns int) *Network {
  tdg := generateTransitDomainGraph(T, Wtt, Wttd)
  return generateTransitDomains(tdg, Nt)
}

// Generate transit domain graph
func generateTransitDomainGraph(T, Wtt, Wttd int) *Network {
  degree := int(math.Log2(float64(T))) + 1
  edges  := degree * T

  mn := Wtt - Wttd
  mx := Wtt + Wttd

  if mn < 0 {
    mn = 0
  }

  return NewRandomUniformNetwork(T, edges, mn, mx)
}

// Generate graph from transit domain
func generateTransitDomains(tdg *Network, Nt int) *Network {
  // Generate the map from node in transit domain graph to
  // graph corresponding to each transit domain
  tdMap := make(map[Router]*Network)

  for _, r := range tdg.Routers {
    nodes := Nt - Ntd + rand.Intn(2 * Ntd + 1)
    if nodes < minNt {
      nodes = minNt
    }

    degree := int(math.Log2(float64(nodes))) + 1
    degree  = degree + rand.Intn(degree * (edgeNtf - 1))
    edges  := degree * nodes

    tdMap[r] = NewRandomUniformNetwork(nodes, edges, minLatency, maxLatency)
  }

  // Generate the new (combined) graph from the two subdomains
  network := new(Network)
  newRouter := make(map[Router]Router)

  for _, nodeNet := range(tdMap) {
    // make map from the node networks to the new combined network
    for _, node := range(nodeNet.Routers) {
      newRouter[node] = NewShortestPathRouter()
      network.Routers = append(network.Routers, newRouter[node])
    }

    // add the edges
    for _, node := range(nodeNet.Routers) {
      for _, conn := range(node.Connections()) {
        n1 := newRouter[node]
        n2 := newRouter[conn.Router()]
        l  := conn.Latency()

        n1.Connect(NewStaticConnection(l, n2))
      }
    }
  }

  // Add the inter-transit edges
  for _, node := range(tdg.Routers) {
    for _, conn := range(node.Connections()) {
      n1 := tdMap[node].RandomRouter()
      n2 := tdMap[conn.Router()].RandomRouter()
      l  := conn.Latency()

      n1.Connect(NewStaticConnection(l, n2))
    }
  }

  return network
}

/* Generates a connected graph.
  Reference: http://economics.mit.edu/files/4622
*/
func NewRandomUniformNetwork(nodes, edges, minLatency, maxLatency int) *Network {
  network := new(Network)

  if math.Log2(float64(nodes)) * float64(nodes) / 2 > float64(edges) {
    panic("Too few number of edges to keep the graph connected.")
  }

  network.Routers = []Router{}
  for i := 0; i < nodes; i++ {
    network.Routers = append(network.Routers, NewShortestPathRouter())
  }

  present := make(map[struct {x, y int}]bool)
  for i := 0; i < edges; i++ {
    i1 := rand.Intn(nodes)
    i2 := rand.Intn(nodes)

    if i1 == i2 {
      i--
      continue
    }

    if present[struct {x, y int}{i1, i2}] || present[struct {x, y int}{i2, i1}] {
      i--
      continue
    } else {
      present[struct {x, y int}{i1, i2}] = true
      present[struct {x, y int}{i2, i1}] = true
    }

    latency := rand.Intn(maxLatency - minLatency) + minLatency
    network.Routers[i1].Connect(NewStaticConnection(latency, network.Routers[i2]))
    network.Routers[i2].Connect(NewStaticConnection(latency, network.Routers[i1]))
  }

  if connected(network) {
    return network
  } else {
    return NewRandomUniformNetwork(nodes, edges, minLatency, maxLatency)
  }
}

func dfs(visited map[Router]bool, router Router) {
  if visited[router] {
    return
  }
  visited[router] = true
  for _, conn := range(router.Connections()) {
    dfs(visited, conn.Router())
  }
}

func connected(net *Network) bool {
  visited := make(map[Router]bool)
  dfs(visited, net.Routers[0])
  ctr := 0
  for _, v := range(visited) {
    if v {
      ctr++
    }
  }
  return ctr == len(net.Routers)
}
