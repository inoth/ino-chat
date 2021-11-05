package service

import (
	"ino-chat/cache"
	"ino-chat/src/chat"
	"ino-chat/util"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func NewRoom(rid, owner, title string) bool {
	err := cache.HMSet(util.RoomInfoKey+rid, map[string]interface{}{
		"rid":   rid,
		"owner": owner,
		"title": title,
	})
	if err != nil {
		log.Errorf("room create err: %v", err)
		return false
	}
	err = cache.SAdd(util.RoomListKey, rid+";"+title)
	if err != nil {
		log.Errorf("save room in list err: %v", err)
		return false
	}
	return JoinRoom(rid, owner)
}

func JoinRoom(rid, uid string) bool {
	if !cache.Exists(util.RoomInfoKey + rid) {
		log.Errorf("room not found.")
		return false
	}
	err := cache.SAdd(util.RoomUserListKey+rid, uid)
	if err != nil {
		log.Errorf("join room err: %v", err)
		return false
	}
	MessageRoomPush(uid, rid, "system", "用户："+uid+"加入了房间")
	return true
}

func ExitRoom(rid, uid string) bool {
	err := cache.SRem(util.RoomUserListKey+rid, uid)
	if err != nil {
		log.Errorf("join room err: %v", err)
		return false
	}
	// 清除用户信息
	ClearUser(uid)
	// 断开连接
	chat.Disconnect(uid)

	// 检查房间内人数
	if !cache.Exists(util.RoomUserListKey + rid) {
		log.Info("No one in the current room.")
		// 删除该房间信息
		removeRoom(rid)
	}
	return true
}

func GetAllRoom() ([]string, error) {
	rooms, err := cache.SMembers(util.RoomListKey)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}
	return rooms, nil
}

func removeRoom(rid string) bool {
	title := cache.HGet(util.RoomInfoKey+rid, "title")
	err := cache.SRem(util.RoomListKey, rid+";"+title)
	if err != nil {
		log.Error(err.Error())
	}
	err = cache.Del(util.RoomInfoKey + rid)
	if err != nil {
		log.Error(err.Error())
	}
	return true
}
