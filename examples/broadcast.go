package examples

import (
	. "github.com/danalex97/Speer/interfaces"
	. "github.com/danalex97/Speer/std"

	"fmt"
)

type BroadcastExample struct {
	Pipeline
	Transport
	Topology

	id string
}

func (s *BroadcastExample) New(util NodeUtil) Node {
	return &BroadcastExample{
		Pipeline: NewChainPipeline(util),
		Transport: util.Transport(),
		Topology: NewA2ATopology(util),

		id: util.Id(),
	}
}

func (s *BroadcastExample) broadcast(m interface{}) {
	for _, member := range s.Neighbors() {
		s.ControlSend(member, m)
	}
}

func (s *BroadcastExample) OnJoin() {
	s.Pipeline = s.Pipeline.
		OnTimeout(s.Topology, 10000).
		Once(func() bool {
			s.broadcast(s.id)
			return true
		}).
		Then(func() bool {
			select {
				case m, _ := <-s.ControlRecv():
					fmt.Println(s.id, "recv:", m)
				default:
			}
			return false
		})
	s.Topology.OnJoin()
}

func (s *BroadcastExample) OnNotify() {
	s.Pipeline.Step()
}

func (s *BroadcastExample) OnLeave() {
}
