package examples

import (
	. "github.com/danalex97/Speer/interfaces"
	. "github.com/danalex97/Speer/std"

	"fmt"
)

type StdMembershipExample struct {
	Pipeline
	Transport
	Membership

	id string
}

func (s *StdMembershipExample) New(util NodeUtil) Node {
	return &StdMembershipExample{
		Pipeline: NewChainPipeline(util),
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
	s.Pipeline = s.Pipeline.
		OnTimeout(s.Membership, 10000).
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
	s.Membership.OnJoin()
}

func (s *StdMembershipExample) OnNotify() {
	s.Pipeline.Step()
}

func (s *StdMembershipExample) OnLeave() {
}
