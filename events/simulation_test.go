package events

import (
  "testing"
)

const offset int = settleTime

func setParallel(t *testing.T, s *Simulation, parallel bool, special bool) {
  if !parallel {
    t.Log("Testing sequential simulation.")
  } else {
    s.check = func(Receiver) bool {
      return special
    }
    t.Log("Testing parallel simulation.")
  }
  s.SetParallel(parallel)
}

func testMultipleTimeReads(t *testing.T, parallel bool, special bool) {
  s := NewLazySimulation()
  setParallel(t, s, parallel, special)

  r := new(mockReceiver)

  go s.Run()

  done := make(chan bool)

  for i := 1; i < LazyQueueChanSize; i++ {
    go func() {
      s.Push(NewEvent(i + offset, nil, r))
      done <- true
      s.Time()
    }()
  }

  for i := 1; i < LazyQueueChanSize; i++ {
    <-done
  }

  s.Stop()
}

func TestMultipleTimeReads(t *testing.T) {
  /*
   * There is no current time ordering guarantee!
   * This is just to make sure there is a locking mechanism over the timer.
   */
  testMultipleTimeReads(t, false, false)
  testMultipleTimeReads(t, true, false)
  testMultipleTimeReads(t, true, true)
}

func testObserversGetNotified(t *testing.T, parallel bool, special bool) {
  r1  := new(mockReceiver)
  r2  := new(mockReceiver)

  o := NewEventObserver(r1)
  s := NewLazySimulation()
  setParallel(t, s, parallel, special)

  go s.Run()
  s.RegisterObserver(o)

  done := make(chan bool)

  go func() {
    for i := 1; i <= LazyQueueChanSize / 2; i++ {
      s.Push(NewEvent(i + offset, nil, r1))
      done <- true
    }
  }()
  go func() {
    for i := 1; i <= LazyQueueChanSize / 2; i++ {
      s.Push(NewEvent(i + offset, nil, r2))
      done <- true
    }
  }()

  for i := 1; i <= LazyQueueChanSize; i++ {
    <-done
  }

  for i := 1; i < LazyQueueChanSize/2; i++ {
    e := (<-o.Recv()).(*Event)
    if e.Timestamp() > LazyQueueChanSize/2 + offset {
      t.Fatalf("Inconsistent simulation times.")
    }
  }

  s.Stop()
}

func TestObserversGetNotified(t *testing.T) {
  /*
   * There is no current time ordering guarantee!
   * Test observers get notified.
   */

  testObserversGetNotified(t, false, false)
  testObserversGetNotified(t, true, false)
  testObserversGetNotified(t, true, true)
}

type oneTimePushReceiver struct {
  pushed bool
}

func (m *oneTimePushReceiver) Receive(e *Event) *Event {
  if !m.pushed {
    m.pushed = true
    return e
  }
  return nil
}

func testReceiversPushNewEvents(t *testing.T, parallel bool, special bool) {
  s := NewLazySimulation()
  setParallel(t, s, parallel, special)

  r := new(oneTimePushReceiver)
  o := NewEventObserver(r)

  go s.Run()
  s.RegisterObserver(o)

  done := make(chan bool)

  for i := 1; i <= LazyQueueChanSize; i++ {
    go func() {
      s.Push(NewEvent(i + offset, nil, r))
      done <- true
    }()
  }

  for i := 1; i <= LazyQueueChanSize; i++ {
    <-done
  }
  for i := 1; i <= LazyQueueChanSize + 1; i++ {
    <-o.Recv()
  }

  s.Stop()
}

func TestReceiversPushNewEvents(t *testing.T) {
  /*
   * There is no current time ordering guarantee!
   * Test receivers push new events if new events get generated.
   */

  testReceiversPushNewEvents(t, false, false)
  testReceiversPushNewEvents(t, true, false)
  testReceiversPushNewEvents(t, true, true)
}

func TestParallelProcessing(t *testing.T) {
  s := NewLazySimulation()
  setParallel(t, s, true, true)

  r := new(oneTimePushReceiver)
  r2 := new(oneTimePushReceiver)
  o := NewEventObserver(r)
  o2 := NewEventObserver(r2)

  go s.Run()
  s.RegisterObserver(o)
  s.RegisterObserver(o2)

  done := make(chan bool)

  for i := 1; i <= LazyQueueChanSize; i++ {
    go func() {
      s.Push(NewEvent(i, nil, r))
      s.Push(NewEvent(i, nil, r2))
      done <- true
    }()
  }

  for i := 1; i <= LazyQueueChanSize; i++ {
    <-done
  }
  for i := 1; i <= LazyQueueChanSize + 1; i++ {
    <-o.Recv()
  }
  for i := 1; i <= LazyQueueChanSize + 1; i++ {
    <-o2.Recv()
  }

  s.Stop()
}
