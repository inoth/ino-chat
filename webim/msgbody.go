package webim

import (
	"encoding/json"
	"errors"
	"time"
	"webchat/webim/msgtype"
)

type MsgBody struct {
	msgType  int
	target   []string
	Uid      string      `json:"from"`
	Timespan int64       `json:"timespan"`
	Body     interface{} `json:"msg"`
}

func (c *MsgBody) ToJson() []byte {
	if by, err := json.Marshal(c); err == nil {
		return by
	}
	return nil
}

type TextBody struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func SendUserMsg(uid, msg string, target ...string) error {
	if len(target) == 0 {
		return errors.New("无效发送对象")
	}
	var t int
	var t1 string
	if len(target) == 1 {
		t = msgtype.MSG
		t1 = "msg"
	} else {
		t = msgtype.MSGS
		t1 = "msgs"
	}
	m := &MsgBody{
		msgType:  t,
		target:   target,
		Uid:      uid,
		Timespan: time.Now().Unix(),
		Body: &TextBody{
			Type: t1,
			Msg:  msg,
		},
	}
	c := Instance()
	c.broadcast <- m
	return nil
}

func SendRoomMsg(uid, msg string, target ...string) error {
	if len(target) == 0 {
		return errors.New("无效房间号")
	}
	var t int
	var t1 string
	if len(target) == 1 {
		t = msgtype.ROOM
		t1 = "room"
	} else {
		t = msgtype.ROOMS
		t1 = "rooms"
	}
	m := &MsgBody{
		msgType:  t,
		target:   target,
		Uid:      uid,
		Timespan: time.Now().Unix(),
		Body: &TextBody{
			Type: t1,
			Msg:  msg,
		},
	}
	c := Instance()
	c.broadcast <- m
	return nil
}
