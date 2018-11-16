package server

import (
  "github.com/danalex97/Speer/events"
  "github.com/danalex97/Speer/overlay"

  "github.com/gorilla/mux"
  "net/http"
  "log"
)

type Server struct {
  monitor *EventMonitor
  router  *mux.Router
}

func NewServer(o events.Observer, netmap *overlay.NetworkMap) *Server {
  return &Server{
    monitor : NewEventMonitor(o, netmap),
    router  : mux.NewRouter(),
  }
}

func (s *Server) Run() {
  go s.monitor.GatherEvents()
  s.router.
    HandleFunc("/new_events", s.monitor.GetNewEvents).
    Methods("Get")
  log.Fatal(http.ListenAndServe(":8000", s.router))
}
