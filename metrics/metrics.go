package metrics

import (
  . "github.com/danalex97/Speer/events"
  "runtime"
  "os"
  "fmt"
)

type Metrics struct {
  events <-chan *Event
}

func NewMetrics(o EventObserver) *Metrics {
  metrics := new(Metrics)
  metrics.events = o.EventChan()
  return metrics
}

func (m *Metrics) Run() {
  f, err := os.OpenFile("metrics.txt", os.O_APPEND|os.O_WRONLY, 0600)
  if err != nil {
    panic(err)
  }
  defer f.Close()

  for {
    select {
    case event := <-m.events:
      line := fmt.Sprintf("%d\n", event.Timestamp())
      if _, err = f.WriteString(line); err != nil {
          panic(err)
      }
    default:
      runtime.Gosched()
    }
  }
}
