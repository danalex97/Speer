package capacity

import (
	"testing"
)

func TestCanConnectCapacityConnectors(t *testing.T) {
	cm := NewScheduledCapacityMap(10)

	cm.AddConnector("1", NewCapacityConnector(10, 20, cm))
	cm.AddConnector("2", NewCapacityConnector(30, 40, cm))

	l := cm.Connector("1").Connect("2")
	assertEqual(t, l.From().Up(), 10)
	assertEqual(t, l.From().Down(), 20)
	assertEqual(t, l.To().Up(), 30)
	assertEqual(t, l.To().Down(), 40)
}
