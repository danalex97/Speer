package std

import (
	. "github.com/danalex97/Speer/interfaces"
)

// Topology primitve. It exposes the local overlay topology visiable to a node. It can 
// be used by a node to work on a subgraph of the original all-to-all overlay topology 
// to implement some algorithms.
type Topology interface {
	Node
	
	Neighbors() []string
}

// Builds full mesh of nodes, exposing the Neighbors() function that can be called by a node.
// Since it reacts all the time to new joins, the simplest way to use it is to set a timeout 
// or use prior knowledge about the number of nodes in the network with a predicate. 
// The implementation uses broadcasts at each new join request arriving at the root of a 
// broadcast tree.
type A2ATopology struct {
	t Transport // we want the transport to be private

	seq    int
	id     string
	parent string

	members []string

	ready   bool
	timeout bool
}

func NewA2ATopology(util NodeUtil) *A2ATopology {
	return &A2ATopology{
		t: util.Transport(),

		seq:    0,
		id:     util.Id(),
		parent: util.Join(),

		members: []string{util.Id()},

		ready:   false,
		timeout: false,
	}
}


func (s *A2ATopology) New(util NodeUtil) Node {
	return NewA2ATopology(util)
}

func (s *A2ATopology) root() bool {
	return s.parent == ""
}

type join struct {
	id string
}

type newMembers struct {
	seq int
	members []string
}

func (s *A2ATopology) Neighbors() []string {
	return s.members
}

func (s *A2ATopology) broadcast(m interface{}) {
	for _, member := range s.members {
		s.t.ControlSend(member, m)
	}
}

func (s *A2ATopology) OnJoin() {
	if !s.root() {
		s.t.ControlSend(s.parent, join{
			id: s.id,
		})
	}
}

func (s *A2ATopology) OnNotify() {
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
				s.seq += 1

				s.broadcast(newMembers{
					members: s.members,
					seq: s.seq,
				})
			}
		case newMembers:
			if !s.root() && msg.seq > s.seq {
				s.members = msg.members
				s.seq = msg.seq
			}
		}
	default:
	}
}

func (s *A2ATopology) OnLeave() {
}
