package sdk

import (
	"github.com/danalex97/Speer/interfaces"

	"runtime"
	"sync"
	"time"
	"os"
	"io/ioutil"
	"strings"

	"testing"
)

var joins int = 0
var messages int = 0
var mutex sync.Mutex = sync.Mutex{}

type mockNode struct {
	id   string
	join string

	transport interfaces.Transport
}

func (s *mockNode) New(util interfaces.NodeUtil) interfaces.Node {
	return &mockNode{
		id:   util.Id(),
		join: util.Join(),

		transport: util.Transport(),
	}
}

func (s *mockNode) OnJoin() {
	mutex.Lock()
	joins += 1
	mutex.Unlock()
	if s.join != "" {
		s.transport.ControlSend(s.join, "hello")
	}

	for {
		select {
		case m, ok := <-s.transport.ControlRecv():
			if !ok {
				continue
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
			runtime.Gosched()
		}
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

	sim := NewSimulationBuilder(new(mockNode)).
		WithInternetworkUnderlay(5, 5, 5, 5).
		WithParallelSimulation().
		WithFixedNodes(10).
		WithCapacityScheduler(1).
		WithCapacityNodes(10, 1, 1).
		Build()

	go sim.Run()
	time.Sleep(100 * time.Millisecond)
	sim.Stop()

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

	go sim.Run()
	time.Sleep(200 * time.Millisecond)
	sim.Stop()

	log, _ := ioutil.ReadFile(file)
	vals := string(log[:])

	if len(strings.Split(vals, "\n")) < 300 {
		t.Fatalf("Log suprisingly short")
	}

	assertEqual(t, joins, 80)
	assertEqual(t, messages, 79)
}
