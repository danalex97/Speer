package metrics

import (
  "github.com/danalex97/Speer/underlay"
  "github.com/danalex97/Speer/overlay"
  "github.com/danalex97/Speer/model"
  . "github.com/danalex97/Speer/events"
  "runtime"
  "os"
  "fmt"
)

type Metrics struct {
  events <-chan *Event
  netmap        *overlay.NetworkMap
}

func NewMetrics(o EventObserver, netmap *overlay.NetworkMap) *Metrics {
  metrics := new(Metrics)
  metrics.events = o.EventChan()
  metrics.netmap = netmap
  return metrics
}

var file = "metrics.txt"

func (m *Metrics) Run() {
  os.Remove(file)
  os.Create(file)

  f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
  if err != nil {
    panic(err)
  }
  defer f.Close()

  for {
    select {
    case event := <-m.events:
      entry := ""

      switch payload := event.Payload().(type) {
      case underlay.Packet:
        underSrc := payload.Src()
        underDst := payload.Dest()
        recv     := m.netmap.Id(event.Receiver().(underlay.Router))
        domain   := event.Receiver().(underlay.Router).Domain()

        src := m.netmap.Id(underSrc)
        dst := m.netmap.Id(underDst)

        // Node ids for the packet are overlay ids.
        entry = fmt.Sprintf("<underlay_packet> src(%s) dest(%s) recv(%s) domain(%s)",
          src, dst, recv, domain)

      case model.DHTQuery:
        key   := payload.Key()
        size  := payload.Size()
        node  := payload.Node()
        store := payload.Store()

        // The given key is a randomly generated key id
        // This SHOULD be changed at the upper layers inside the protocol
        // implmentation
        entry = fmt.Sprintf("<query> key(%s) size(%d) node(%s) store(%t)",
          key, size, node, store)

      case model.Join:
        nodeId     := payload.NodeId()

        entry = fmt.Sprintf("<join> nodeId(%s)", nodeId)
      default:
        entry = "<nil>"
      }

      if entry == "<nil>" {
        // ignore other entries
        continue
      }

      // timestamping the entry
      entry = fmt.Sprintf("%d %s\n", event.Timestamp(), entry)

      if _, err = f.WriteString(entry); err != nil {
          panic(err)
      }

    default:
      runtime.Gosched()
    }
  }
}
