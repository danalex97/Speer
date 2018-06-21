package underlay

import (
  "container/list"
  "errors"
)

type Cache interface {
  Put(key, value interface {})
  Get(key interface {}) (interface {}, bool)
}

// A Cache using the Least Recently Used policy. It is implemented by using a
// list and a map pointing to list elements.
type LRUCache struct {
  size  int
  evict *list.List
  items map[interface{}]*list.Element
}

type centry struct {
  key   interface{}
  value interface{}
}

func NewLRUCache(size int) (*LRUCache, error) {
  if size <= 0 {
    return nil, errors.New("Non-positive cache size")
  }

  c := &LRUCache {
    size:  size,
    evict: list.New(),
    items: make(map[interface{}]*list.Element),
  }
  return c, nil
}

func (c *LRUCache) Put(key, value interface{}) {
  // Look for present element
  if entry, ok := c.items[key]; ok {
    c.evict.MoveToFront(entry)
    entry.Value.(*centry).value = value
    return
  }

  // Add new element
  entry     := &centry{key, value}
  listEntry := c.evict.PushFront(entry)
  c.items[key] = listEntry

  // Evict element
  if c.evict.Len() > c.size {
    // Remove from list
    listEntry := c.evict.Back()
    c.evict.Remove(listEntry)

    // Remove from map
    entry := listEntry.Value.(*centry)
    delete(c.items, entry.key)
  }
}

func (c *LRUCache) Get(key interface{}) (interface{}, bool) {
  if entry, ok := c.items[key]; ok {
    c.evict.MoveToFront(entry)
    return entry.Value.(*centry).value, true
  }
  return nil, false
}
