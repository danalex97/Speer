package examples

import (
	. "github.com/danalex97/Speer/interfaces"
	str "github.com/danalex97/Speer/structs"

	"fmt"
)

type DataLinkExample struct {
	Transport

	node1 string
	node2 string

	time func() int
}

func (s *DataLinkExample) New(util NodeUtil) Node {
	return &DataLinkExample{
		Transport: util.Transport(),

		node1: util.Id(),
		node2: util.Join(),

		time: util.Time(),
	}
}

func (s *DataLinkExample) OnJoin() {
	if s.node2 != "" {
		s.ControlSend(s.node2, s.node1)
		s.ControlSend(s.node1, s.node2)
	}

	var up Link
	var down Link

	for {
		switch msg := (<-s.ControlRecv()).(type) {
		case string:
			// establish uplink when I know ID of other person
			up = s.Connect(msg)
			s.ControlSend(msg, up)
		case Link:
			// establish downlink when I receive ACK
			down = msg
		}

		if up != nil && down != nil {
			fmt.Println("Node", s.node1, "capacities", s.Up(), s.Down())

			// Upload some pieces of data
			size := 10000
			for i := 1; i <= 5; i++ {
				up.Upload(Data{str.RandomKey()[:3], size})
			}

			// Download some pieces of data
			for i := 1; i <= 5; i++ {
				// Print the data I received at a timestamp
				fmt.Println("Node", s.node1, ":", <-down.Download())
			}
		}
	}
}

func (s *DataLinkExample) OnLeave() {
}
