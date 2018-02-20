package model

import (
  "testing"
  "math"
  "math/rand"
  "time"
)

type genFunc func() float64

func testGeneratorPoissonMean(t *testing.T, g genFunc, rate float64) {
  rand.Seed(time.Now().UTC().UnixNano())

  samples := 1000
  eps := 0.1

  avg := 0.0
  for i := 0; i < samples; i++ {
    avg += g()
  }
  avg = float64(avg) / float64(samples)
  if math.Max(avg - rate, rate - avg) / rate > eps {
    t.Fatalf("Poisson rate does not equal sample mean: %s != %s", rate, avg)
  }
}

func TestPoissonProcessModelNextArrivalCloseToRate(t *testing.T) {
  for _, arrRate := range []float64{40.0, 100.0, 500.0, 30.0} {
    testGeneratorPoissonMean(t,
      NewPoissonProcessModel(arrRate, 0.0).NextArrival, arrRate)
  }
}

func TestPoissonProcessModelNextDepartureCloseToRate(t *testing.T) {
  for _, arrRate := range []float64{40.0, 100.0, 500.0, 30.0} {
    testGeneratorPoissonMean(t,
      NewPoissonProcessModel(arrRate, 0.0).NextDeparture, arrRate)
  }
}

func TestPoissonProcessModelNextQueryCloseToRate(t *testing.T) {
  for _, queryRate := range []float64{40.0, 100.0, 500.0, 30.0} {
    testGeneratorPoissonMean(t,
      NewPoissonProcessModel(0.0, queryRate).NextQuery, queryRate)
  }
}
