package logs

import (
  . "github.com/danalex97/Speer/events"

  "github.com/danalex97/Speer/underlay"
  "github.com/danalex97/Speer/overlay"
  "github.com/danalex97/Speer/model"

  "encoding/json"
  "runtime"
  "fmt"
  "os"
)

const eventQueueCapacity = 1000000
const maxEvents = 100

type EventMonitor struct {
  newEvents chan interface{}

  incomingEvents <-chan interface{}
  netmap *overlay.NetworkMap
  outFile string
}

func NewEventMonitor(
    o Observer,
    netmap *overlay.NetworkMap,
    outFile string,
  ) *EventMonitor {
  return &EventMonitor{
    newEvents  : make(chan interface{}, eventQueueCapacity),

    incomingEvents : o.Recv(),
    netmap         : netmap,
    outFile        : outFile,
  }
}

func (em *EventMonitor) GatherEvents() {
  os.Remove(em.outFile)
  os.Create(em.outFile)

  f, err := os.OpenFile(em.outFile, os.O_APPEND|os.O_WRONLY, 0644)
  if err != nil {
    panic(err)
  }
  defer f.Close()

  for {
    select {
    case msg := <-em.incomingEvents:
      event := msg.(*Event)
      timestamp := event.Timestamp()

      var newEvent interface{}

      switch payload := event.Payload().(type) {
      case underlay.Packet:
        underSrc := payload.Src()
        underDst := payload.Dest()
        recv     := em.netmap.Id(event.Receiver().(underlay.Router))
        router   := event.Receiver().(underlay.Router)

        src := em.netmap.Id(underSrc)
        dst := em.netmap.Id(underDst)

        newEvent = UnderlayPacketEntry{
          Time   : timestamp,

          Src : src,
          Dst : dst,
          Rtr : recv,

          SrcUid : fmt.Sprintf("%p", (&underSrc)),
          DstUid : fmt.Sprintf("%p", (&underDst)),
          RtrUid : fmt.Sprintf("%p", (&router)),
        }
      case model.Join:
        nodeId := payload.NodeId()

        newEvent = JoinEntry{
          Time : timestamp,
          Node : nodeId,
        }
      }

      if (newEvent != nil) {
        son, _ := json.Marshal(newEvent)
        f.WriteString(string(son))
        f.WriteString("\n")
      }
    default:
      runtime.Gosched()
    }
  }
}