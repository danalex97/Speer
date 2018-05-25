package events

import (
  "testing"
)

func setParallel(t *testing.T, s Simulation, parallel bool) {
  if parallel {
    t.Log("Testing sequential simulation.")
  } else {
    t.Log("Testing parallel simulation.")
  }
  s.SetParallel(parallel)
}

func TestMultipleTimeReads(t *testing.T) {
  /*
   * There is no current time ordering guarantee!
   * This is just to make sure there is a locking mechanism over the timer.
   */
  test := func(parallel bool) {
    s := NewLazySimulation()
    setParallel(t, s, parallel)

    r := new(mockReceiver)

    go s.Run()

    done := make(chan bool)

    for i := 1; i < LazyQueueChanSize; i++ {
      go func() {
        s.Push(NewEvent(i, nil, r))
        done <- true
        s.Time()
      }()
    }

    for i := 1; i < LazyQueueChanSize; i++ {
      <-done
    }

    s.Stop()
  }

  test(false)
  test(true)
}

func TestObserversGetNotified(t *testing.T) {
  /*
   * There is no current time ordering guarantee!
   * Test observers get notified.
   */

  test := func(parallel bool) {
    r1  := new(mockReceiver)
    r2  := new(mockReceiver)

    o := NewEventObserver(r1)
    s := NewLazySimulation()
    setParallel(t, s, parallel)

    go s.Run()
    s.RegisterObserver(o)

    done := make(chan bool)

    go func() {
      for i := 1; i <= LazyQueueChanSize / 2; i++ {
        s.Push(NewEvent(i, nil, r1))
        done <- true
      }
    }()
    go func() {
      for i := 1; i <= LazyQueueChanSize / 2; i++ {
        s.Push(NewEvent(i, nil, r2))
        done <- true
      }
    }()

    for i := 1; i <= LazyQueueChanSize; i++ {
      <-done
    }

    for i := 1; i < LazyQueueChanSize/2; i++ {
      e := (<-o.Recv()).(*Event)
      if e.Timestamp() > LazyQueueChanSize/2 {
    		t.Fatalf("Inconsistent simulation times.")
    	}
      // For assserting time ordering
      // assertEqual(t, e.Timestamp(), i)
    }

    s.Stop()
  }

  test(false)
  test(true)
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

func TestReceiversPushNewEvents(t *testing.T) {
  /*
   * There is no current time ordering guarantee!
   * Test receivers push new events if new events get generated.
   */

   test := func(parallel bool) {
     s := NewLazySimulation()
     setParallel(t, s, parallel)

     r := new(oneTimePushReceiver)
     o := NewEventObserver(r)

     go s.Run()
     s.RegisterObserver(o)

     done := make(chan bool)

     for i := 1; i <= LazyQueueChanSize; i++ {
       go func() {
         s.Push(NewEvent(i, nil, r))
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

  test(false)
  test(true)
}
