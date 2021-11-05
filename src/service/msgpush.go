package service

import (
	"ino-chat/cache"
	"ino-chat/src/chat"
	"ino-chat/src/models"
	"ino-chat/util"

	"github.com/sirupsen/logrus"
)

// 单对单发送消息
func MessagePush(uid, target, msgType, msg string) {
	body := &models.MsgPushBody{
		Target:  []string{target},
		FromUid: uid,
		MsgType: msgType,
		Body:    msg,
	}
	chat.SendMsg(body)
}

// 根据传入房间信息和发送人，给整个房间内人推送消息
func MessageRoomPush(uid, rid, msgType, msg string) {
	body := &models.MsgPushBody{
		Target:  GetUsersByRoom(rid),
		FromUid: uid,
		MsgType: msgType,
		Body:    msg,
	}
	chat.SendMsg(body)
}

func GetUsersByRoom(rid string) []string {
	users, err := cache.SMembers(util.RoomUserListKey + rid)
	if err != nil {
		logrus.Error()
		return []string{}
	}
	return users
}
