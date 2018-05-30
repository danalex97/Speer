package structs

import (
  "testing"
)

func assertNotEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Fatalf("%s == %s", a, b)
	}
}

func TestRandomKeyUniqueKeys(t *testing.T) {
  keys := []string{}
  for i := 0; i < 100; i++ {
    keys = append(keys, RandomKey())
  }

  for i := 0; i < 100; i++ {
    for j := 0; j < i; j++ {
      if keys[i] == keys[j] {
        assertNotEqual(t, keys[i], keys[j])
      }
    }
  }
}
