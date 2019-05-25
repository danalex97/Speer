package examples

import (
	. "github.com/danalex97/Speer/interfaces"

	"fmt"
)

type RoutinesExample struct {
	RoutineCapabilities
	Transport

	node1 string
	node2 string

	time func() int
}

func (s *RoutinesExample) New(util NodeUtil) Node {
	return &RoutinesExample{
		RoutineCapabilities: util,
		Transport: util.Transport(),

		node1: util.Id(),
		node2: util.Join(),

		time: util.Time(),
	}
}

func (s *RoutinesExample) OnJoin() {
	if s.node2 != "" {
		t := 0
		r := s.Routine(10, func() {
			s.ControlSend(s.node2, t)
			t += 1
		})
		s.Callback(100, func() {
			r.SetInterval(20)
		})
		s.Callback(200, func() {
			r.Stop()
		})
	}
}

func (s *RoutinesExample) OnNotify() {
	select {
	case m, _ := <-s.ControlRecv():
		fmt.Println("@", s.time(), "received:", m)
	default:
	}
}

func (s *RoutinesExample) OnLeave() {
}
