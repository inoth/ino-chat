package webim

import "encoding/json"

const (
	TEXT = iota
	IMG
	VOICE
	VIDIO
)

type MsgBody struct {
	msgType   int32 // 0:system  1:user
	target    []string
	Uid       string `json:"from"`
	Timestamp int64  `json:"timestamp"`
	Body      string `json:"msg"`
}

func (c *MsgBody) ToJson() []byte {
	if by, err := json.Marshal(c); err == nil {
		return by
	}
	return nil
}

func SendMessage(msgType int32, timestamp int64, target []string, fromUser, body string) error {
	Instance().broadcast <- &MsgBody{
		msgType:   msgType,
		target:    target,
		Timestamp: timestamp,
		Uid:       fromUser,
		Body:      body,
	}
	return nil
}
