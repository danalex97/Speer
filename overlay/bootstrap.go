package overlay

import (
  "strconv"
  "math/rand"
  . "github.com/danalex97/Speer/underlay"
)

type Bootstrap interface {
  Join(id string) string // returns bootstrap node
  NewId() string // returns node ID
  Router(id string) Router // ID -> Router
  Id(router Router) string // Router -> ID
}

type NetworkMap struct {
  network *Network
  id map[string]Router
  inv map[Router]string
  idCtr int
}

func NewNetworkMap(network *Network) Bootstrap {
  mp := new(NetworkMap)

  mp.network = network
  mp.id      = make(map[string]Router)
  mp.inv     = make(map[Router]string)
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

func (mp *NetworkMap) Router(id string) Router {
  if router, ok := mp.id[id]; ok{
    return router
  } else {
    return nil
  }
}

func (mp *NetworkMap) Id(router Router) string {
  if id, ok := mp.inv[router]; ok{
    return id
  } else {
    return ""
  }
}
