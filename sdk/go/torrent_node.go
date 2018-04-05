package sdk

import (
  "github.com/danalex97/Speer/capacity"
)

type TorrentNode interface {
  DHTNode

  Transfer() capacity.Engine
}

type AutowiredTorrentNode struct {
  AutowiredDHTNode
}

func (n *AutowiredDHTNode) Transfer() capacity.Engine {
  return nil
}
