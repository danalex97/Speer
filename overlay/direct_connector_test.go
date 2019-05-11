package overlay

import (
	"testing"
)

func testDirectConnector() (string, string, DirectConnector, DirectConnector) {
	chanMap := NewChanMap()

	bridge1, id1 := NewDirectChan(chanMap)
	bridge2, id2 := NewDirectChan(chanMap)

	return id1, id2, bridge1, bridge2
}

func TestDirectConnector(t *testing.T) {
	_, id2, bridge1, bridge2 := testDirectConnector()

	for i := 0; i < 10; i++ {
		t.Logf("DirectChann packet delivery test -- packet %d\n", i)
		bridge1.ControlSend(id2, "message")
		assertEqual(t, "message", <-bridge2.ControlRecv())
	}
}

func TestDirectConnectorSendPacketToSelf(t *testing.T) {
	for i := 0; i < 10; i++ {
		id1, _, bridge1, _ := testDirectConnector()

		bridge1.ControlSend(id1, "message")
		assertEqual(t, "message", <-bridge1.ControlRecv())
	}
}
