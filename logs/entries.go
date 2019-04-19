package logs

type UnderlayPacketEntry struct {
	Time int `json:"time"`

	Src string `json:"src"`
	Dst string `json:"dst"`
	Rtr string `json:"rtr"`

	SrcUid string `json:"src_uid"`
	DstUid string `json:"dst_uid"`
	RtrUid string `json:"rtr_uid"`
}

type JoinEntry struct {
	Time int    `json:"time"`
	Node string `json:"node"`
}
