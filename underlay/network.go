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

func NewRandomUniformNetwork(nodes, edges, minLatency, maxLatency int) *Network {
  // reference: http://economics.mit.edu/files/4622
  network := new(Network)

  if math.Log2(float64(nodes)) * float64(nodes) / 2 > float64(edges) {
    panic("Too few number of edges to keep the graph connex.")
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
