package overlay

import (
  "strconv"
  "math/rand"
  . "github.com/danalex97/Speer/underlay"
)

type Bootstrap interface {
  Join() string // returns bootstrap node
  Id() string // returns node ID
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

func (mp *NetworkMap) Id() string {
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
