package model

type MessageNsqBody struct {
	MsgType  int      `json:"msgType"`
	Target   []string `json:"target"`
	FromUser string   `json:"fromUser"`
	Body     string   `json:"body"`
}
