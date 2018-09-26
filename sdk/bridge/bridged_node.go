package bridge

import (
  . "github.com/danalex97/Speer/interfaces"
)

type BridgedTorrent struct {
  envChannel <-chan Message
  bridge     EnvironBridge

  TorrentNodeUtil
}

func (t *BridgedTorrent) New(util TorrentNodeUtil) TorrentNode {
  return &BridgedTorrent{
    TorrentNodeUtil : util,

    envChannel      : util.Bridge().RecvMessage(util.Id()),
    bridge          : util.Bridge(),
  }
}

func (t *BridgedTorrent) OnJoin() {
  t.bridge.SendMessage(&Create{
    Id : t.Id(),
  })
}

func (t *BridgedTorrent) OnLeave() {
}
