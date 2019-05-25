package examples

import (
	.   "github.com/danalex97/Speer/interfaces"
	str "github.com/danalex97/Speer/structs"

	"fmt"
)

type DataLinkExample struct {
	Transport

	node1 string
	node2 string

	up Link
	down Link

	time func() int
}

func (s *DataLinkExample) New(util NodeUtil) Node {
	return &DataLinkExample{
		Transport : util.Transport(),

		node1 : util.Id(),
		node2 : util.Join(),

		time : util.Time(),
	}
}

func (s *DataLinkExample) OnJoin() {
	if s.node2 != "" {
		// node2 will initiate both connections
		s.ControlSend(s.node2, s.node1)
		s.ControlSend(s.node1, s.node2)
	}
}

func (s *DataLinkExample) OnNotify() {
	select {
	case m, _ := <-s.ControlRecv():
		switch msg := m.(type) {
		case string:
			// establish uplink when I know ID of other person
			s.up = s.Connect(msg)
			s.ControlSend(msg, s.up)
		case Link:
			// establish downlink when I receive ACK
			s.down = msg
		default:
		}
	default:
	}

	// after I connected, upload the data
	if s.up != nil && s.down != nil {
		fmt.Println("Node", s.node1, "capacities", s.Up(), s.Down())

		// Upload some pieces of data
		size := 100
		for i := 1; i <= 5; i++ {
			s.up.Upload(Data{str.RandomKey()[:3], size})
		}

		// stop the upload
		s.up = nil;
	}

	// receive the downloaded data
	if s.down != nil {
		select {
		case data, _ := <-s.down.Download():
			fmt.Println("@", s.time(), "node", s.node1, ":", data)
		default:
		}
	}
}

func (s *DataLinkExample) OnLeave() {
}
