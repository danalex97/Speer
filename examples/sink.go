package examples

import (
	. "github.com/danalex97/Speer/interfaces"
	
	"runtime"
	"fmt"
)

type SinkExample struct {
	Transport

	id string
	parent string
}

func (s *SinkExample) New(util NodeUtil) Node {
	return &SinkExample{
		Transport : util.Transport(),

		id : util.Id(),
		parent : util.Join(),
	}
}

func (s *SinkExample) root() bool {
	return s.parent == ""
}

func (s *SinkExample) handleRecv(m interface{}) {
	if !s.root() {
		// forward each message
		s.ControlSend(s.parent, m)
	} else {
		// the root will print the messages
		fmt.Println("Received", m)
	}
}

func (s *SinkExample) OnJoin() {
	// send my id to the parent
	if !s.root() {
		s.ControlSend(s.parent, s.id)
	}

	for {
		select {
		case m, _ := <-s.ControlRecv():
			s.handleRecv(m);

		default:
			runtime.Gosched()
		}
	}
}

func (s *SinkExample) OnLeave() {
}
