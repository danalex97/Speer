package metrics

import (
  "github.com/danalex97/Speer/underlay"
  // "github.com/danalex97/Speer/overlay"
  . "github.com/danalex97/Speer/events"
  "runtime"
  "os"
  "fmt"
)

type Metrics struct {
  events <-chan *Event
  // bridge
}

func NewMetrics(o EventObserver) *Metrics {
  metrics := new(Metrics)
  metrics.events = o.EventChan()
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

        entry = fmt.Sprintf("packet: %s %s\n", underSrc, underDst)
      }

      if _, err = f.WriteString(entry); err != nil {
          panic(err)
      }
    default:
      runtime.Gosched()
    }
  }
}
