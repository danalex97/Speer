package sdk

import (
  "github.com/danalex97/Speer/capacity"
)

type TorrentNode interface {
  DHTNode

  Transfer() capacity.Engine

  // used to autowire the engine
  autowireEngine(capacity.Engine)
}

type AutowiredTorrentNode struct {
  AutowiredDHTNode

  engine capacity.Engine
}

func (a *AutowiredTorrentNode) Transfer() capacity.Engine {
  return a.engine
}

func (a *AutowiredTorrentNode) autowireEngine(engine capacity.Engine) {
  a.engine = engine
}
