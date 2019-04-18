package underlay

import (
	"testing"
)

func TestNewPacket(t *testing.T) {
	r := NewShortestPathRouter("1")
	p := NewPacket(r, r, nil)

	assertEqual(t, p.Src(), r)
	assertEqual(t, p.Dest(), r)
	assertEqual(t, p.Payload(), nil)
}
