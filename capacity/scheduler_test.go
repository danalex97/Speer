package capacity

import (
  . "github.com/danalex97/Speer/interfaces"
  "testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
    t.Fatalf("%s != %s", a, b)
	}
}

/* Mock structures. */

type link struct {
  from int
  to   int
}

/* Scheduler full scenario test. */

func TestSchedulerFullScenario(t *testing.T) {
  s := NewScheduler(10).(*scheduler)

  nodes := []Node{
    &node{10, 0},
    &node{0, 30},
    &node{0, 20},
  }
  links := []Link{
    NewPerfectLink(nodes[1], nodes[0]),
    NewPerfectLink(nodes[2], nodes[0]),
  }
  s.RegisterLink(links[0])
  s.RegisterLink(links[1])

  // time: 0
  links[0].Upload(Data{"1-0-0", 100})
  links[1].Upload(Data{"2-0-0", 50})
  s.Schedule()

  // time: 10
  links[1].Upload(Data{"2-0-1", 50})
  s.Schedule()
  if len(links[0].Download()) != 0 {
    t.Fatalf("Link 1-0 transmitted: %s", len(links[0].Download()))
  }
  if len(links[1].Download()) != 1 {
    t.Fatalf("Link 2-0 not transmitted: %s", len(links[1].Download()))
  }

  // time: 20
  links[0].Upload(Data{"1-0-1", 100})
  s.Schedule()
  if len(links[0].Download()) != 1 {
    t.Fatalf("Link 1-0 not transmitted: %s", len(links[0].Download()))
  }
  if len(links[1].Download()) != 2 {
    t.Fatalf("Link 2-0 not transmitted: %s", len(links[1].Download()))
  }

  // time: 30
  s.Schedule()
  if len(links[0].Download()) != 2 {
    t.Fatalf("Link 1-0 not transmitted: %s", len(links[0].Download()))
  }
  if len(links[1].Download()) != 2 {
    t.Fatalf("Link 2-0 not transmitted: %s", len(links[1].Download()))
  }

  assertEqual(t, Data{"1-0-0", 100}, <-links[0].Download())
  assertEqual(t, Data{"1-0-1", 100}, <-links[0].Download())
  assertEqual(t, Data{"2-0-0", 50}, <-links[1].Download())
  assertEqual(t, Data{"2-0-1", 50}, <-links[1].Download())
}

/* Capacity upload checks. */

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

func TestUpdCapacityMultipleToOne(t *testing.T) {
  buildTest(t, []node{
    node{5, 0},
    node{0, 2},
    node{0, 2},
  }, []link{
    link{1, 0},
    link{2, 0},
  }, func(s *scheduler, nodes []node, links []Link) {
    s.updCapacity()
    checkCapacity(t, s, links[0], 2)
    checkCapacity(t, s, links[1], 2)
  })

  buildTest(t, []node{
    node{5, 0},
    node{0, 3},
    node{0, 3},
  }, []link{
    link{1, 0},
    link{2, 0},
  }, func(s *scheduler, nodes []node, links []Link) {
    s.updCapacity()
    checkCapacity(t, s, links[0], 2.5)
    checkCapacity(t, s, links[1], 2.5)
  })

  buildTest(t, []node{
    node{5, 0},
    node{0, 10},
    node{0, 10},
  }, []link{
    link{1, 0},
    link{2, 0},
  }, func(s *scheduler, nodes []node, links []Link) {
    s.updCapacity()
    checkCapacity(t, s, links[0], 2.5)
    checkCapacity(t, s, links[1], 2.5)
  })
}
