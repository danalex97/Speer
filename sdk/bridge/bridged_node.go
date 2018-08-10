package bridge

import (
  . "github.com/danalex97/Speer/interfaces"
)

func getEnviron() *Environ {
  return nil
}

type BridgedTorrent struct {
  envChannel chan<- interface {}

  TorrentNodeUtil
}

func (t *BridgedTorrent) New(util TorrentNodeUtil) TorrentNode {
  return &BridgedTorrent{
    TorrentNodeUtil : util,
    envChannel      : make(chan interface {}),
  }
}

func (t *BridgedTorrent) OnJoin() {
}

func (t *BridgedTorrent) OnLeave() {
}
