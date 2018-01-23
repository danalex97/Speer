package overlay

import (
  "strconv"
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

func newId(mp *NetworkMap) (id string) {
  mp.idCtr++
  id = strconv.Itoa(mp.idCtr)
  return
}

func (mp *NetworkMap) Id() string {
  for {
    if router, ok := inv[mp.network.RandomRouter()] ; ok {
      routerId := newId(mp)

      mp.id[routerId] = router
      mp.inv[router]  = routerId

      return routerId
    }
  }
}
