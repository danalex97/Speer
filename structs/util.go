package structs

import (
  "math/rand"
  "runtime"
)

func RandomKey() string {
  const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

  b := make([]byte, 30)
  for i := range b {
    b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
  }
  return string(b)
}

func Wait(cond func () bool) {
  for {
    if cond() {
      runtime.Gosched()
    } else {
      break
    }
  }
}
