package webim

import "encoding/json"

type MsgBody struct {
	msgType   int
	target    []int
	Uid       string `json:"from"`
	Timestamp int32  `json:"timestamp"`
	Body      string `json:"msg"`
}

func (c *MsgBody) ToJson() []byte {
	if by, err := json.Marshal(c); err == nil {
		return by
	}
	return nil
}
