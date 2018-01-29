package model

import (
  "math"
  "math/rand"
)

type DHTQueryGenerator interface {
  Next() *DHTQuery
  // the key for a store query is empty, thus
  // allowing the SDK layer to handle it
}

// each query should be immutable
type DHTQuery struct {
  key   string // the key of the node
  size  int    // size of key to be transfered in MB
  node  string // the node which sends/stores the query
  store bool   // store/retrieve
}

func NewDHTQuery(key string, size int, node string, store bool) *DHTQuery{
  q := new(DHTQuery)
  q.key = key
  q.size = size
  q.node = node
  q.store = store
  return q
}

func (q *DHTQuery) Key() string {
  return q.key
}

func (q *DHTQuery) Size() int {
  return q.size
}

func (q *DHTQuery) Node() string {
  return q.node
}

func (q *DHTQuery) Store() bool {
  return q.store
}

type DHTLedger struct {
  queries   []*DHTQuery
  bootstrap Bootstrap
}

func NewDHTLedger(bootstrap Bootstrap) *DHTLedger {
  ledger := new(DHTLedger)
  ledger.queries = []*DHTQuery{}
  ledger.bootstrap = bootstrap
  return ledger
}

func (l *DHTLedger) Next() *DHTQuery {
  node := bootstrap.Join()
  size := rand.Int31n(100)
  store := length(queries) == 0 || rand.Float32() > 0.5
  key   := ""

  if !store {
    // this is generated uniformly as there are no leaves yet
    // and the history has only 'store' queries
    idx := rand.Int31n(length(l.queries))
    key = queries[idx].Key()
  }

  query := NewDHTQuery(key, size, node, store)
  if store {
    l.queries = append(queries, query)
  }

  return query
}
