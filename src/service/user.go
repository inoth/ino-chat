package service

import (
	"ino-chat/cache"
	"ino-chat/util"

	log "github.com/sirupsen/logrus"
)

func SaveUserInCache(uid string) bool {
	err := cache.HMSet(util.UserInfoKey+uid, map[string]interface{}{
		"uid": uid,
	}, 60*60*24)
	if err != nil {
		log.Errorf("login err: %v", err)
		return false
	}
	return true
}

func ClearUser(uid string) bool {
	err := cache.Del(util.UserInfoKey + uid)
	if err != nil {
		log.Errorf("login err: %v", err)
		return false
	}
	return true
}
