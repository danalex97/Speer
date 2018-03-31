package underlay

import (
  "math"
  "math/rand"
)

type Network struct {
  Routers []Router
  stubRouters []Router
}

func (n *Network) RandomRouter() Router {
  if n.stubRouters == nil {
    return n.Routers[rand.Intn(len(n.Routers))]
  } else {
    return n.stubRouters[rand.Intn(len(n.stubRouters))]
  }
}

/* Constants used in the stub-generation algorithm.

  Wtt  = avg. weight of transit-transit connections
  Wttd = Wtt delta, that is the weight is in [Wtt - Wttd, Wtt + Wttd]

  Ntd     = nodes transit domain delta (see Wttd)
  minNt   = minimum number of nodes in transit domain
  edgeNtf = edge factor for a transit domain, that is the number of edges
    is f * N log N

  minLatency = minimum latency for intra-domain connections
  maxLatency = maximum latency for intra-domain connections

  Nsd     = ...
  minNs   = ...
  edgeNsf = ...

  mhsp = percent of multi-home stub *connections*

  The constant values are currently arbirary, but they can be chosen to respect
  the invariants in the paper.
*/
const Wtt  int = 100
const Wttd int = 2

const Ntd     int = 5
const minNt   int = 5
const edgeNtf int = 2

const minLatency int = 2
const maxLatency int = 10

const Wts  int = 100
const Wtsd int = 2

const Nsd     int = 5
const minNs   int = 5
const edgeNsf int = 2

const mhsp int = 50

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
  backbone := generateTransitDomains(tdg, Nt)
  network := addStubs(backbone, S, Ns)
  return addMhs(network, S)
}

// 1. Generate transit domain graph
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

// 2. Generate graph from transit domain
func generateTransitDomains(tdg *Network, Nt int) *Network {
  // Generate the map from node in transit domain graph to
  // graph corresponding to each transit domain
  tdMap := make(map[Router]*Network)

  for _, r := range tdg.Routers {
    tdMap[r] = newRandomNetwork(Nt, minNt, Ntd, edgeNtf, minLatency, maxLatency)
  }

  // Generate the new (combined) graph from the two subdomains
  network := new(Network)
  newRouter := make(map[Router]Router)

  for _, nodeNet := range tdMap {
    copyNetwork(newRouter, network, nodeNet)
  }

  // Add the inter-transit edges
  present := make(map[struct {x, y Router}]bool)
  for _, node := range tdg.Routers {
    for _, conn := range node.Connections() {
      n1 := tdMap[node].RandomRouter()
      n2 := tdMap[conn.Router()].RandomRouter()
      l  := conn.Latency()

      // We need this as we want to add bidirectional edges
      if !present[struct {x, y Router}{n1, n2}] {
        newRouter[n1].Connect(NewStaticConnection(l, newRouter[n2]))
        newRouter[n2].Connect(NewStaticConnection(l, newRouter[n1]))
      }

      present[struct {x, y Router}{n1, n2}] = true
      present[struct {x, y Router}{n2, n1}] = true
    }
  }

  return network
}

// 3. Add stubs
func addStubs(backbone *Network, S, Ns int) *Network {
  // copy backbone
  network := new(Network)
  for _, node := range backbone.Routers {
    network.Routers = append(network.Routers, node)
  }

  newRouter := make(map[Router]Router)
  for i := 0; i < S; i++ {
    // generate stub
    stub := newRandomNetwork(Ns, minNs, Nsd, edgeNsf, minLatency, maxLatency)

    // copy stub on network
    copyNetwork(newRouter, network, stub)

    // attach stub to network
    attachBack := backbone.RandomRouter()
    attachStub := newRouter[stub.RandomRouter()]

    addTsEdge(attachBack, attachStub)
  }

  // make list of stub nodes
  stubRouters := []Router{}
  for _, router := range newRouter {
    stubRouters = append(stubRouters, router)
  }
  network.stubRouters = stubRouters

  return network
}

// 4. Add multi-homed stubs
func addMhs(network *Network, stubs int) *Network {
  mhs := mhsp * stubs / 100
  stubNodes := network.stubRouters

  stubSet := make(map[Router]bool)
  for _, node := range stubNodes {
    stubSet[node] = true
  }
  backNodes := []Router{}
  for _, node := range network.Routers {
    if _, ok := stubSet[node]; !ok {
      backNodes = append(backNodes, node)
    }
  }

  // Add mhs random edges from a stub node to a backbone node
  for i := 0; i < mhs; i++ {
    stubNode := stubNodes[rand.Intn(len(stubNodes))]
    backNode := backNodes[rand.Intn(len(backNodes))]

    addTsEdge(backNode, stubNode)
  }

  return network
}

/* Helper functions */
func newRandomNetwork(nodes, minNodes, nodesDelta, edgeFactor, minLatency, maxLatency int) *Network {
  nodes = nodes - nodesDelta + rand.Intn(2 * nodesDelta + 1)
  if nodes < minNodes {
    nodes = minNodes
  }

  degree := int(math.Log2(float64(nodes))) + 1
  degree  = degree + rand.Intn(degree * (edgeFactor - 1))
  edges  := degree * nodes

  return NewRandomUniformNetwork(nodes, edges, minLatency, maxLatency)
}

// Adds toCopy connex component to network and updates the mapping newRouter
func copyNetwork(newRouter map[Router]Router, network *Network, toCopy *Network) {
  // make map from the node networks to the new combined network
  for _, node := range toCopy.Routers {
    newRouter[node] = NewShortestPathRouter()
    network.Routers = append(network.Routers, newRouter[node])
  }

  // add the edges
  for _, node := range toCopy.Routers {
    for _, conn := range node.Connections() {
      n1 := newRouter[node]
      n2 := newRouter[conn.Router()]
      l  := conn.Latency()

      n1.Connect(NewStaticConnection(l, n2))
    }
  }
}

func addTsEdge(attachBack, attachStub Router) {
  minLatency := Wts - Wtsd
  maxLatency := Wts + Wtsd
  latency := rand.Intn(maxLatency - minLatency) + minLatency

  attachStub.Connect(NewStaticConnection(latency, attachBack))
  attachBack.Connect(NewStaticConnection(latency, attachStub))
}

const smallNetNodes int = 20

func insertEdge(present map[struct {x, y int}]bool, network *Network, i1, i2, minL, maxL int) bool {
  if i1 == i2 {
    return false
  }

  if present[struct {x, y int}{i1, i2}] || present[struct {x, y int}{i2, i1}] {
    return false
  } else {
    present[struct {x, y int}{i1, i2}] = true
    present[struct {x, y int}{i2, i1}] = true
  }

  latency := rand.Intn(maxL - minL) + minL
  network.Routers[i1].Connect(NewStaticConnection(latency, network.Routers[i2]))
  network.Routers[i2].Connect(NewStaticConnection(latency, network.Routers[i1]))

  return true
}

func smallRandomNetwork(nodes, edges, minLatency, maxLatency int) *Network {
  network := new(Network)
  present := make(map[struct {x, y int}]bool)

  network.Routers = []Router{}
  for i := 0; i < nodes; i++ {
    network.Routers = append(network.Routers, NewShortestPathRouter())
  }

  // build tree
  for i := 1; i < nodes; i++ {
    i1 := rand.Intn(i)
    i2 := i

    insertEdge(present, network, i1, i2, minLatency, maxLatency)
  }

  // add rest of edges
  for i := 0; i < edges; i++ {
    i1 := rand.Intn(nodes)
    i2 := rand.Intn(nodes)

    // since we have only a few nodes, we might not add all the edges
    insertEdge(present, network, i1, i2, minLatency, maxLatency)
  }

  return network
}

/* Generates a connected graph.
  Reference: http://economics.mit.edu/files/4622
*/
func NewRandomUniformNetwork(nodes, edges, minLatency, maxLatency int) *Network {
  if math.Log2(float64(nodes)) * float64(nodes) / 2 > float64(edges) {
    panic("Too few number of edges to keep the graph connected.")
  }

  if (nodes < smallNetNodes) {
    return smallRandomNetwork(nodes, edges, minLatency, maxLatency)
  }

  network := new(Network)
  network.Routers = []Router{}
  for i := 0; i < nodes; i++ {
    network.Routers = append(network.Routers, NewShortestPathRouter())
  }

  present := make(map[struct {x, y int}]bool)
  for i := 0; i < edges; i++ {
    i1 := rand.Intn(nodes)
    i2 := rand.Intn(nodes)

    if !insertEdge(present, network, i1, i2, minLatency, maxLatency) {
      i--
    }
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
  for _, conn := range router.Connections() {
    dfs(visited, conn.Router())
  }
}

func connected(net *Network) bool {
  visited := make(map[Router]bool)
  dfs(visited, net.Routers[0])
  ctr := 0
  for _, v := range visited {
    if v {
      ctr++
    }
  }
  return ctr == len(net.Routers)
}
