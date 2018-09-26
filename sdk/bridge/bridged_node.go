package bridge

import (
  . "github.com/danalex97/Speer/interfaces"
  "fmt"
)

func getEnviron() *Environ {
  return nil
}

type BridgedTorrent struct {
  envChannel chan<- interface {}

  TorrentNodeUtil
}

func (t *BridgedTorrent) New(util TorrentNodeUtil) TorrentNode {
  fmt.Println("new")
  return &BridgedTorrent{
    TorrentNodeUtil : util,
    envChannel      : make(chan interface {}),
  }
}

func (t *BridgedTorrent) OnJoin() {
  fmt.Println("join")
}

func (t *BridgedTorrent) OnLeave() {
}
