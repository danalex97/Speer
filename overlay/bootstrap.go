package overlay

import (
  "strconv"
  "math/rand"
  "github.com/danalex97/Speer/underlay"
)

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
  network *underlay.Network
  id map[string]underlay.Router
  inv map[underlay.Router]string
  idCtr int
}

func NewNetworkMap(network *underlay.Network) OverlayMap {
  mp := new(NetworkMap)

  mp.network = network
  mp.id      = make(map[string]underlay.Router)
  mp.inv     = make(map[underlay.Router]string)
  mp.idCtr   = 0

  return mp
}

func newId(mp *NetworkMap) (id string) {
  mp.idCtr++
  id = strconv.Itoa(mp.idCtr)
  return
}

func (mp *NetworkMap) NewId() string {
  for {
    router := mp.network.RandomRouter()
    if _, ok := mp.inv[router]; !ok {
      routerId := newId(mp)

      mp.id[routerId] = router
      mp.inv[router]  = routerId

      return routerId
    }
  }
}

func (mp *NetworkMap) Join(id string) string {
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
  panic("Join method called on invalid mp.")
}

func (mp *NetworkMap) Router(id string) underlay.Router {
  if router, ok := mp.id[id]; ok{
    return router
  } else {
    return nil
  }
}

func (mp *NetworkMap) Id(router underlay.Router) string {
  if id, ok := mp.inv[router]; ok{
    return id
  } else {
    return ""
  }
}
