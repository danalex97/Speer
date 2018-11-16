package server

import (
  "encoding/json"
  "net/http"
)

func GetTestEvents(w http.ResponseWriter, r *http.Request) {
  events := []interface{}{JoinEntry{
    Time : 0,
    Node : "1",
  }, JoinEntry{
    Time : 1,
    Node : "3",
  }, UnderlayPacketEntry{
    Time : 2,
    Src    : "1",     Dst  : "3",        Rtr : "1",
    SrcUid : "0x1", DstUid : "0x3", RtrUid  : "0x1",
  }, UnderlayPacketEntry{
    Time : 3,
    Src    : "1",     Dst  : "3",        Rtr : "2",
    SrcUid : "0x1", DstUid : "0x3", RtrUid  : "0x2",
  }, UnderlayPacketEntry{
    Time : 4,
    Src    : "1",     Dst  : "3",        Rtr : "3",
    SrcUid : "0x1", DstUid : "0x3", RtrUid  : "0x3",
  }}

  json.NewEncoder(w).Encode(events)
}
