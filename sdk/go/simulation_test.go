package sdk

import (
	"github.com/danalex97/Speer/interfaces"

	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"testing"
)

var joins int = 0
var messages int = 0
var ctr int = 0
var mutex sync.Mutex = sync.Mutex{}

type mockNode struct {
	id   string
	join string

	callback interfaces.Callback
	routine  interfaces.Routine

	transport interfaces.Transport
}

func (s *mockNode) New(util interfaces.NodeUtil) interfaces.Node {
	r := &mockNode{
		id:   util.Id(),
		join: util.Join(),

		callback: nil,
		routine:  nil,

		transport: util.Transport(),
	}
	r.callback = util.Callback(5, func() {
		mutex.Lock()
		ctr += 1
		mutex.Unlock()
	})
	r.routine = util.Routine(5, func() {
		mutex.Lock()
		ctr += 1
		mutex.Unlock()
		r.routine.Stop()
	})
	return r
}

func (s *mockNode) OnJoin() {
	mutex.Lock()
	joins += 1
	mutex.Unlock()
	if s.join != "" {
		s.transport.ControlSend(s.join, "hello")
	}
}

func (s *mockNode) OnNotify() {
	select {
	case m, ok := <-s.transport.ControlRecv():
		if !ok {
			return
		}

		switch msg := m.(type) {
		case string:
			if msg == "hello" {
				mutex.Lock()
				messages += 1
				mutex.Unlock()
			}
		}
	default:
	}
}

func (s *mockNode) OnLeave() {
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestSimulationBuilderAndTransports(t *testing.T) {
	joins = 0
	messages = 0
	ctr = 0

	sim := NewSimulationBuilder(new(mockNode)).
		WithInternetworkUnderlay(5, 5, 5, 5).
		WithParallelSimulation().
		WithFixedNodes(10).
		WithCapacityScheduler(1).
		WithCapacityNodes(10, 1, 1).
		Build()

	sim.Run()
	time.Sleep(200 * time.Millisecond)
	sim.Stop()

	assertEqual(t, ctr, 20)
	assertEqual(t, joins, 10)
	assertEqual(t, messages, 9)
}

func TestSimulationNoTopology(t *testing.T) {
	joins = 0
	messages = 0
	ctr = 0

	sim := NewSimulationBuilder(new(mockNode)).
		WithParallelSimulation().
		WithFixedNodes(10).
		WithCapacityScheduler(1).
		WithCapacityNodes(10, 1, 1).
		Build()

	sim.Run()
	time.Sleep(200 * time.Millisecond)
	sim.Stop()

	assertEqual(t, ctr, 20)
	assertEqual(t, joins, 10)
	assertEqual(t, messages, 9)
}

func TestSimulationNoCapacities(t *testing.T) {
	joins = 0
	messages = 0
	ctr = 0

	sim := NewSimulationBuilder(new(mockNode)).
		WithParallelSimulation().
		WithFixedNodes(10).
		Build()

	sim.Run()
	time.Sleep(200 * time.Millisecond)
	sim.Stop()

	assertEqual(t, ctr, 20)
	assertEqual(t, joins, 10)
	assertEqual(t, messages, 9)
}

func TestSimulationOnFlatToplogy(t *testing.T) {
	file := "log.txt"

	defer func() {
		os.Remove(file)
	}()

	joins = 0
	messages = 0
	ctr = 0

	sim := NewSimulationBuilder(new(mockNode)).
		WithRandomUniformUnderlay(200, 1000, 5, 10).
		WithParallelSimulation().
		WithFixedNodes(100).
		WithCapacityScheduler(1).
		WithCapacityNodes(10, 1, 1).
		WithCapacityNodes(20, 1, 1).
		WithCapacityNodes(50, 1, 1).
		WithLogs(file).
		Build()

	sim.Run()
	time.Sleep(200 * time.Millisecond)
	sim.Stop()

	log, _ := ioutil.ReadFile(file)
	vals := string(log[:])

	assertEqual(t, ctr, 160)
	assertEqual(t, joins, 80)
	assertEqual(t, messages, 79)

	if len(strings.Split(vals, "\n")) < 200 {
		t.Fatalf("Log suprisingly short")
	}
}
