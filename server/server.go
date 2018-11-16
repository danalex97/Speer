package server

import (
  "github.com/danalex97/Speer/events"
  "github.com/danalex97/Speer/overlay"

  "github.com/gorilla/mux"
  "net/http"
  "log"

  "fmt"
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

func (s *Server) Run() {
  go s.monitor.GatherEvents()

  s.router.
    HandleFunc("/new_events", s.monitor.GetNewEvents).
    Methods("Get")

  path, _ := os.Getwd()
  fmt.Println(path)

  path += STATIC_DIR
  dr := http.Dir(path)

  fmt.Println(dr)

  s.router.
    PathPrefix(STATIC_ROUTE).
    Handler(http.StripPrefix(STATIC_ROUTE,
      http.FileServer(dr)))

  log.Fatal(http.ListenAndServe(":" + PORT, s.router))
}
