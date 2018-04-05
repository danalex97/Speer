package examples

import (
  "math/rand"
  "runtime"
)

func Wait(cond func () bool) {
  for {
    if cond() {
      runtime.Gosched()
    } else {
      break
    }
  }
}

func RandomKey() string {
  const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

  b := make([]byte, 30)
  for i := range b {
    b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
  }
  return string(b)
}
