package underlay

import (
  "testing"
)

func TestLRUCacheCanGetAndPutElements(t *testing.T) {
  c, err := NewLRUCache(128)
  if err != nil {
    t.Fatalf("err: %v", err)
  }

  for i := 0; i < 50; i++ {
    c.Put(i, i)
  }

  for i := 0; i < 50; i++ {
    el, ok := c.Get(i)
    assertEqual(t, el, i)
    assertEqual(t, ok, true)
  }
}

func TestLRUOldElementsGetEvicted(t *testing.T) {
  c, err := NewLRUCache(128)
  if err != nil {
    t.Fatalf("err: %v", err)
  }

  for i := 0; i < 256; i++ {
    c.Put(i, i)
  }

  for i := 0; i < 128; i++ {
    _, ok := c.Get(i)
    assertEqual(t, ok, false)
  }
  for i := 128; i < 256; i++ {
    el, ok := c.Get(i)
    assertEqual(t, el, i)
    assertEqual(t, ok, true)
  }
}

func TestLRUCacheAccessedElementsDontGetEvicted(t *testing.T) {
  c, err := NewLRUCache(128)
  if err != nil {
    t.Fatalf("err: %v", err)
  }

  for i := 0; i < 128; i++ {
    c.Put(i, i)
  }
  for i := 0; i < 64; i++ {
    c.Get(i)
  }
  for i := 128; i < 192; i++ {
    c.Put(i, i)
  }

  for i := 64; i < 128; i++ {
    _, ok := c.Get(i)
    assertEqual(t, ok, false)
  }
  for i := 0; i < 64; i++ {
    el, ok := c.Get(i)
    assertEqual(t, el, i)
    assertEqual(t, ok, true)
  }
  for i := 128; i < 192; i++ {
    el, ok := c.Get(i)
    assertEqual(t, el, i)
    assertEqual(t, ok, true)
  }
}
