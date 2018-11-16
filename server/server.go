package server

import (
  "github.com/danalex97/Speer/events"
  "github.com/danalex97/Speer/overlay"

  "github.com/gorilla/mux"
  "net/http"
  "log"

  "os"
)

const (
  STATIC_ROUTE = "/display"
  STATIC_DIR   = "/server/static/"
  PORT         = "8000"
)

type Server struct {
  monitor *EventMonitor
  router  *mux.Router
}

func NewServer(o events.Observer, netmap *overlay.NetworkMap) *Server {
  return &Server{
    monitor : NewEventMonitor(o, netmap),
    router  : mux.NewRouter().StrictSlash(true),
  }
}

func NewTestServer() *Server {
  return &Server{
    router : mux.NewRouter().StrictSlash(true),
  }
}

func (s *Server) bindStatic() {
  path, _ := os.Getwd()
  path    += STATIC_DIR
  dr      := http.Dir(path)

  s.router.
    PathPrefix(STATIC_ROUTE).
    Handler(http.StripPrefix(STATIC_ROUTE,
      http.FileServer(dr)))
}

func (s *Server) TestRun() {
  s.router.
    HandleFunc("/new_events", GetTestEvents).
    Methods("Get")
  s.bindStatic()

  log.Fatal(http.ListenAndServe(":" + PORT, s.router))
}

func (s *Server) Run() {
  go s.monitor.GatherEvents()

  s.router.
    HandleFunc("/new_events", s.monitor.GetNewEvents).
    Methods("Get")
  s.bindStatic()

  log.Fatal(http.ListenAndServe(":" + PORT, s.router))
}
