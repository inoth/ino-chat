package models

import "ino-chat/util"

type LoginData struct {
	Uid    string `json:"-" redis:"uid"`
	Name   string `json:"name" redis:"name"`
	Avatar string `json:"avatar" redis:"avatar"`
}

type NewRoomBody struct {
	Title string `json:"title"`
}

type RoomInfo struct {
	Rid   string `json:"rid"`
	Owner string `json:"owner"`
	Title string `json:"title"`
}

type MsgSendBody struct {
	Target     string `json:"target"`
	TargetType string `json:"targetType"`
	Msg        string `json:"msg"`
	MsgType    string `json:"msgType"`
}

type MsgPushBody struct {
	Target  []string `json:"target"`
	FromUid string   `json:"fromUid"`
	MsgType string   `json:"msgType"`
	Body    string   `json:"body"`
}

type WsSendMessage struct {
	FromUid string `json:"fromUid"`
	MsgType string `json:"msgType"`
	Body    string `json:"body"`
}

func (m *WsSendMessage) ToJson() []byte {
	return util.ToJson(m)
}
