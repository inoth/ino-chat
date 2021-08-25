package model

type RoomInfo struct {
	Rid   string `json:"rid"`
	RName string `json:"rname"`
}

type MessageNsqBody struct {
	MsgType    int    `json:"msgType"`
	TargetType int    `json:"targetType"`
	Target     string `json:"target"`
	FromUser   string `json:"fromUser"`
	Body       string `json:"body"`
}
