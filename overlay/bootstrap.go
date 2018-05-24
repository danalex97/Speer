package overlay

import (
  "github.com/danalex97/Speer/events"
  "github.com/danalex97/Speer/underlay"
  "github.com/danalex97/Speer/interfaces"

  "sync"
  "strconv"
  "math/rand"
)

type GroupProgress interfaces.GroupProgress

// The delay for pushing a packet into the network is at most 100ms.
const progressInterval int = 100

type Bootstrap interface {
  Join(id string) string // returns bootstrap node
}

type OverlayMap interface {
  Bootstrap

  NewId() string // returns node ID

  Router(id string) underlay.Router // ID -> underlay.Router
  Id(router underlay.Router) string // underlay.Router -> ID
}

type NetworkMap struct {
  sync.Mutex

  network *underlay.Network

  id    map[string]underlay.Router
  inv   map[underlay.Router]string

  idCtr map[string]int

  // Progress properties that regulate the transfer of packets.
  pushProgress GroupProgress
  pullProgress GroupProgress
}

func NewNetworkMap(network *underlay.Network) OverlayMap {
  mp := new(NetworkMap)

  mp.network = network
  mp.id      = make(map[string]underlay.Router)
  mp.inv     = make(map[underlay.Router]string)
  mp.idCtr   = make(map[string]int)

  mp.pushProgress = events.NewWGProgress(progressInterval)
  mp.pullProgress = events.NewWGProgress(progressInterval)

  return mp
}

func newId(mp *NetworkMap, domain string) (id string) {
  if _, ok := mp.idCtr[domain]; !ok {
    mp.idCtr[domain] = 0
  }

  id = strconv.Itoa(mp.idCtr[domain])
  if domain != "" {
    id = domain + "." + id
  }

  mp.idCtr[domain]++

  return
}

func (mp *NetworkMap) NewId() string {
  mp.Lock()
  defer mp.Unlock()

  for {
    router := mp.network.RandomRouter()
    if _, ok := mp.inv[router]; !ok {
      routerId := newId(mp, router.Domain())

      mp.id[routerId] = router
      mp.inv[router]  = routerId

      return routerId
    }
  }
}

func (mp *NetworkMap) Join(id string) string {
  // mp.Lock()
  // defer mp.Unlock()

  i := rand.Intn(len(mp.id))
  for k := range(mp.id) {
    if i == 0 {
      if k == id && len(mp.id) > 1 {
        return mp.Join(id)
      }
      if k != id {
        return k
      }
    }
    i--
  }
  return ""
}

func (mp *NetworkMap) Router(id string) underlay.Router {
  mp.Lock()
  defer mp.Unlock()

  if router, ok := mp.id[id]; ok{
    return router
  } else {
    return nil
  }
}

func (mp *NetworkMap) Id(router underlay.Router) string {
  mp.Lock()
  defer mp.Unlock()

  if id, ok := mp.inv[router]; ok {
    return id
  } else {
    return ""
  }
}
