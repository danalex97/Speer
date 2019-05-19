package examples

import (
	. "github.com/danalex97/Speer/interfaces"

	"fmt"
	"runtime"
)

type BroadcastExample struct {
	Transport

	id     string
	parent string

	members []string
	time    func() int
}

func (s *BroadcastExample) New(util NodeUtil) Node {
	return &BroadcastExample{
		Transport: util.Transport(),

		id:     util.Id(),
		parent: util.Join(),

		members: []string{util.Id()},
		time:    util.Time(),
	}
}

func (s *BroadcastExample) root() bool {
	return s.parent == ""
}

type Join struct {
	id string
}

type NewMember struct {
	id string
}

type SomeBroadcast struct {
	ts   int
	list []string
	from string
}

func (s *BroadcastExample) broadcast(m interface{}) {
	for _, member := range s.members {
		s.ControlSend(member, m)
	}
}

func (s *BroadcastExample) handleRecv(m interface{}) {
	switch msg := m.(type) {
	case Join:
		if !s.root() {
			// subscribe in the list of nodes
			s.ControlSend(s.parent, msg)
		} else {
			// if the root receives a new node, broadcast the message
			s.members = append(s.members, msg.id)
			s.broadcast(NewMember{
				id: msg.id,
			})
		}
	case NewMember:
		if !s.root() {
			if msg.id != s.id {
				s.members = append(s.members, msg.id)
				s.broadcast(SomeBroadcast{
					ts:   s.time(),
					list: s.members,
					from: s.id,
				})
			}
		}
	case SomeBroadcast:
		fmt.Println(s.id, "recv:", msg)
	}
}

func (s *BroadcastExample) OnJoin() {
	if !s.root() {
		s.ControlSend(s.parent, Join{
			id: s.id,
		})
	}

	for {
		select {
		case m, _ := <-s.ControlRecv():
			s.handleRecv(m)

		default:
			runtime.Gosched()
		}
	}
}

func (s *BroadcastExample) OnLeave() {
}
