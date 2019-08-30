package examples

import (
	. "github.com/danalex97/Speer/interfaces"
	. "github.com/danalex97/Speer/std"

	"fmt"
)

type StdMembershipExample struct {
	Transport
	Membership

	id string
	once bool
}

func (s *StdMembershipExample) New(util NodeUtil) Node {
	return &StdMembershipExample{
		Transport: util.Transport(),
		Membership: NewBroadcastMembership(util),
	
		id: util.Id(),
		once: false,
	}
}

func (s *StdMembershipExample) broadcast(m interface{}) {
	for _, member := range s.Members() {
		s.ControlSend(member, m)
	}
}

func (s *StdMembershipExample) OnJoin() {
	s.Membership.OnJoin()
}

func (s *StdMembershipExample) OnNotify() {
	if s.Membership.ComposeOnTimeout(10000) {
		if !s.once {
			s.broadcast(s.id)
			s.once = true
		}

		select {
		case m, _ := <-s.ControlRecv():
			fmt.Println(s.id, "recv:", m)
		default:
		}
	}
}

func (s *StdMembershipExample) OnLeave() {
}
