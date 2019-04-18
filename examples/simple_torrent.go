package examples

import (
	. "github.com/danalex97/Speer/interfaces"
	"github.com/danalex97/Speer/structs"
	"runtime"
	"sync"
	"fmt"
)

type SimpleTorrent struct {
	sync.Mutex

	id         string
	ids        []string
	links      map[string]Link
	transport  Transport
}

type idBroadcast struct {
	ids []string
}

/* Interface functions. */
func (s *SimpleTorrent) OnJoin() {
	go func() {
		for {
			// check links
			for _, l := range s.links {
				select {
				case m, _ := <-l.Download():
					fmt.Println(s.id, "data", m)
				default:
					continue
				}
			}

			// check for control messages
			select {

			case m, ok := <-s.transport.ControlRecv():
				if !ok {
					continue
				}

				switch msg := m.(type) {
				case idBroadcast:
					s.updateIds(msg.ids)
					fmt.Println(s.id, "received", msg.ids)
				}

			default:
				runtime.Gosched()
			}
		}
	}()

	go func() {
		// broadcast neighbours
		for _, id := range s.ids {
			if id != s.id {
				if !s.transport.ControlPing(id) {
					continue
				}

				s.transport.ControlSend(id, idBroadcast{s.ids})
			}
		}
	}()
}

func (s *SimpleTorrent) OnLeave() {
}

func (s *SimpleTorrent) New(util NodeUtil) Node {
	// Constructor that assumes the UnreliableNode component is filled in
	node := new(SimpleTorrent)

	node.id = util.Id()
	node.ids = []string{node.id}
	join := util.Join()
	if join != "" {
		node.ids = append(node.ids, join)
	}

	node.transport = util.Transport()
	node.links = map[string]Link{}

	return node
}

/* Local functions */
func (s *SimpleTorrent) updateIds(ids []string) {
	allIds := make(map[string]bool)
	for _, id := range ids {
		allIds[id] = true
	}
	for _, id := range s.ids {
		allIds[id] = true
	}

	s.ids = []string{}
	for id, _ := range allIds {
		s.ids = append(s.ids, id)
		
		if id == s.id {
			continue
		}

		// register link if not registered
		if _, ok := s.links[id]; !ok {
			s.links[id] = s.transport.Connect(id)

			// if the link is new, we broadcast our list again
			if !s.transport.ControlPing(id) {
				continue
			}
			s.transport.ControlSend(id, idBroadcast{s.ids})

			// send a big packet
			s.links[id].Upload(Data{structs.RandomKey(), 1000})
		}
	}
}
