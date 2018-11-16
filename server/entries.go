package server

type UnderlayPacketEntry struct {
  Time   int `json:"time"`

  Src string `json:"src"`
  Dst string `json:"dst"`
  Rtr string `json:"rtr"`

  SrcUid int32 `json:"rtr_uid"`
  DstUid int32 `json:"rtr_uid"`
  RtrUid int32 `json:"rtr_uid"`
}

type JoinEntry struct {
  Time int `json:"time"`
  Node string `json:"node"`
}
