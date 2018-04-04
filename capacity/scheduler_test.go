package capacity

import (
  "testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

type node struct {
  down int
  up   int
}

type link struct {
  from int
  to   int
}

func (n *node) Up() int {
  return n.up
}

func (n *node) Down() int {
  return n.down
}

func buildTest(t *testing.T, nodes []node, idxs []link, callback func(*scheduler, []node, []Link)) {
  s := NewScheduler(0).(*scheduler)

  links := []Link{}
  for _, l := range idxs {
    link := NewPerfectLink(&nodes[l.from], &nodes[l.to])
    links = append(links, link)
    s.RegisterLink(link)
  }
  for _, status := range s.linkStatus {
    status.active = true
  }

  callback(s, nodes, links)
}

func checkCapacity(t *testing.T, s *scheduler, link Link, cap float64) {
  if s.linkStatus[link] == nil {
    t.Fatalf("Link not found!")
  }
  assertEqual(t, s.linkStatus[link].capacity, cap)
}

func TestUpdCapacityTwoNodes(t *testing.T) {
  buildTest(t, []node{
    node{10, 10},
    node{10, 10},
  }, []link{
    link{0, 1},
  }, func(s *scheduler, nodes []node, links []Link) {
    s.updCapacity()
    checkCapacity(t, s, links[0], 10)
  })

  buildTest(t, []node{
    node{10, 2},
    node{10, 10},
  }, []link{
    link{0, 1},
  }, func(s *scheduler, nodes []node, links []Link) {
    s.updCapacity()
    checkCapacity(t, s, links[0], 2)
  })

  buildTest(t, []node{
    node{10, 10},
    node{2,  10},
  }, []link{
    link{0, 1},
  }, func(s *scheduler, nodes []node, links []Link) {
    s.updCapacity()
    checkCapacity(t, s, links[0], 2)
  })
}

func TestUpdCapacitySimpleGraph(t *testing.T) {
  /**
  (down | up)

  (0 | 10) --->  (8 | 0)
      |              ^
      |              |
      +-->(3 | 2) ---+
   */

  buildTest(t, []node{
    node{0, 10},
    node{3, 2},
    node{8, 0},
  }, []link{
    link{0, 1},
    link{1, 2},
    link{0, 2},
  }, func(s *scheduler, nodes []node, links []Link) {
    s.updCapacity()
    checkCapacity(t, s, links[0], 3)
    checkCapacity(t, s, links[1], 2)
    checkCapacity(t, s, links[2], 6)
  })
}
