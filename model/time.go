package model

type TimeModel interface {
  NextArrival() int    // next time a node arrival happens
  NextDeparture() int  // next time a node departure happends
  NextQuery() int      // next time a query happens
}
