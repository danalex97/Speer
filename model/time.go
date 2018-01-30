package model

import (
  "math"
  "math/rand"
)

type TimeModel interface {
  NextArrival() float64    // next time a node arrival happens
  NextDeparture() float64  // next time a node departure happends
  NextQuery() float64      // next time a query happens
}

type PoissonProcessModel struct {
  arrivalRate float64 // arrivals/second
  queryRate float64   // queries/second
}

func NewPoissonProcessModel(arrivalRate float64, queryRate float64) *PoissonProcessModel {
  p := new(PoissonProcessModel)
  p.arrivalRate = arrivalRate
  p.queryRate   = queryRate
  return p
}

func poissonNext(rate float64) float64 {
  return -math.Log(1.0 - rand.Float64()) * rate;
}

func (p *PoissonProcessModel) NextArrival() float64 {
  return poissonNext(p.arrivalRate)
}

func (p *PoissonProcessModel) NextDeparture() float64 {
  return poissonNext(p.arrivalRate)
}

func (p *PoissonProcessModel) NextQuery() float64 {
  return poissonNext(p.queryRate)
}
