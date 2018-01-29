package model

import (
  "math"
  "math/rand"
)

type TimeModel interface {
  NextArrival() int    // next time a node arrival happens
  NextDeparture() int  // next time a node departure happends
  NextQuery() int      // next time a query happens
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

func poissonNext(rate float64) {
  return -math.Log(1.0 - rand.float64()) / rate;
}

func (p *PoissonProcessModel) NextArrival() {
  return poissonNext(p.arrivalRate)
}

func (p *PoissonProcessModel) NextDeparture() {
  return poissonNext(p.arrivalRate)
}

func (p *PoissonProcessModel) NextQuery() {
  return poissonNext(p.queryRate)
}
