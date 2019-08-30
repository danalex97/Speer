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
}

func (s *StdMembershipExample) New(util NodeUtil) Node {
	return &StdMembershipExample{
		Transport: util.Transport(),
		Membership: NewBroadcastMembership(util),
	
		id: util.Id(),
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
	s.Membership.ComposeOnTimeout(100)
	s.broadcast(s.id)

	select {
	case m, _ := <-s.ControlRecv():
		fmt.Println(s.id, "recv:", m)
	default:
	}
}

func (s *StdMembershipExample) OnLeave() {
}
