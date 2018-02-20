package model

import (
  "testing"
)

func assertEqual(t *testing.T, a, b interface {}) {
  if a != b {
    t.Fatalf("%s == %s", a, b)
  }
}

func assertNotEqual(t *testing.T, a, b interface {}) {
  if a == b {
    t.Fatalf("%s == %s", a, b)
  }
}

type mockBootstrap struct {}
func (m *mockBootstrap) Join(id string) string {
  return ""
}

func TestDHTLedgerEqualQueryTypes(t *testing.T) {
  tot := 1000
  per := 400

  storeQ  := 0
  deleteQ := 0
  gen     := NewDHTLedger(new(mockBootstrap))

  for i := 0; i < tot; i++ {
    query := gen.Next()
    if query.Store() {
      storeQ += 1
    } else {
      deleteQ += 1
    }
  }

  if storeQ < per {
    t.Fatalf("Too few store queries: %s of %s", storeQ, tot)
  }
  if deleteQ < per {
    t.Fatalf("Too few delete queries: %s of %s", deleteQ, tot)
  }
}

func TestDHTLedgerDeleteQueriesFromPerviouslySeenKey(t *testing.T) {
  // Node we have no guarantee on size, and should not be used
  for i := 1; i < 100; i++ {
    gen := NewDHTLedger(new(mockBootstrap))
    query  := gen.Next()
    query2 := gen.Next()
    if query2.Store() {
      assertNotEqual(t, query.Key(), query2.Key())
    } else {
      assertEqual(t, query.Key(), query2.Key())
    }
  }
}

func TestDHTLedgerQuerySizeDoesNotExceedMax(t *testing.T) {
  gen := NewDHTLedger(new(mockBootstrap))
  for i := 1; i < 100; i++ {
    query := gen.Next()
    if query.Size() > MaxQuerySize {
      t.Fatalf("Query size exceeded: %s", query.Size())
    }
  }
}
