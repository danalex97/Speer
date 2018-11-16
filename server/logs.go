package server

import (
  . "github.com/danalex97/Speer/events"

  "github.com/danalex97/Speer/underlay"
  "github.com/danalex97/Speer/overlay"
  "github.com/danalex97/Speer/model"

  "encoding/json"
  "net/http"
  "runtime"
  "unsafe"
)

type EventMonitor struct {
  newEvents chan interface{}

  incomingEvents <-chan interface{}
  netmap *overlay.NetworkMap
}

func NewEventMonitor(o Observer, netmap *overlay.NetworkMap) *EventMonitor {
  return &EventMonitor{
    newEvents  : make(chan interface{}),

    incomingEvents : o.Recv(),
    netmap         : netmap,
  }
}

func (em *EventMonitor) GatherEvents() {
  for {
    select {
    case msg := <-em.incomingEvents:
      event := msg.(*Event)
      timestamp := event.Timestamp()

      switch payload := event.Payload().(type) {
      case underlay.Packet:
        underSrc := payload.Src()
        underDst := payload.Dest()
        recv     := em.netmap.Id(event.Receiver().(underlay.Router))
        router   := event.Receiver().(underlay.Router)

        src := em.netmap.Id(underSrc)
        dst := em.netmap.Id(underDst)

        em.newEvents <- UnderlayPacketEntry{
          Time   : timestamp,

          Src : src,
          Dst : dst,
          Rtr : recv,

          SrcUid : *((*int32)(unsafe.Pointer(&underSrc))),
          DstUid : *((*int32)(unsafe.Pointer(&underDst))),
          RtrUid : *((*int32)(unsafe.Pointer(&router))),
        }
      case model.Join:
        nodeId := payload.NodeId()

        em.newEvents <- JoinEntry{
          Time : timestamp,
          Node : nodeId,
        }
      }
    default:
      runtime.Gosched()
    }
  }
}

func (em *EventMonitor) GetNewEvents(w http.ResponseWriter, r *http.Request) {
  packets := []UnderlayPacketEntry{}
  json.NewEncoder(w).Encode(packets)
}
