package underlay

import (
  "testing"
)

func TestStaticConnectionLatencyRemainsConstant(t *testing.T) {
  r := NewShortestPathRouter("1")
  l := 10
  c := NewStaticConnection(l, r)

  assertEqual(t, c.Router(), r)
  for i := 0; i < 10; i++ {
    assertEqual(t, c.Latency(), l)
  }
}
