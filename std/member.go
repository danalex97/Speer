package std

import (
	. "github.com/danalex97/Speer/interfaces"
)

type Composable interface {
	ComposeOnNotify(timeout int)
}

// Membership primitve. Builds full mesh of nodes, exposing the Members() function that can 
// be called by a node. Implements the node interface. Since it reacts all the time to new 
// joins, the simplest way to use it is to set a timeout or use prior knowledge about the 
// number of nodes in the network. 
type Membership interface {
	Composable
	Node
	
	Members() []string
}

// Membership primitive implementation using broadcasts at each new join request arriving 
// at the root of a broadcast tree.
type BroadcastMembership struct {
	r RoutineCapabilities
	t Transport // we want the transport to be private

	id     string
	parent string

	members []string

	ready   bool
	timeout bool
}

func (s *BroadcastMembership) New(util NodeUtil) Node {
	return &BroadcastMembership{
		t: util.Transport(),

		id:     util.Id(),
		parent: util.Join(),

		members: []string{util.Id()},

		ready: false,
		timeout: false,
	}
}

func (s *BroadcastMembership) root() bool {
	return s.parent == ""
}

type join struct {
	id string
}

type newMembers struct {
	members []string
}

func (s *BroadcastMembership) Members() []string {
	return s.members
}

func (s *BroadcastMembership) broadcast(m interface{}) {
	for _, member := range s.members {
		s.t.ControlSend(member, m)
	}
}

func (s *BroadcastMembership) OnJoin() {
	if !s.root() {
		s.t.ControlSend(s.parent, join{
			id: s.id,
		})
	}

}

// Composition is done by putting a basic timeout. 
// [TODO:] Separate in own class -- timeout composable 
func (s *BroadcastMembership) ComposeOnNotify(timeout int) {
	if s.ready {
		return
	}

	if !s.timeout {
		s.timeout = true
		s.r.Callback(timeout, func() {
			s.ready = true
		})
	}

	s.OnNotify()
}

func (s *BroadcastMembership) OnNotify() {
	select {
	case m, _ := <-s.t.ControlRecv():
		switch msg := m.(type) {
		case join:
			if !s.root() {
				// subscribe in the list of nodes
				s.t.ControlSend(s.parent, msg)
			} else {
				// if the root receives a new node, broadcast the message
				s.members = append(s.members, msg.id)
				s.broadcast(newMembers{
					members: s.members,
				})
				// fmt.Println("Broadcast")
			}
		case newMembers:
			if !s.root() {
				s.members = msg.members
				// fmt.Println(s.id, s.members)
			}
		}
	default:
	}
}

func (s *BroadcastMembership) OnLeave() {
}
