package util

import (
	"encoding/json"

	"github.com/google/uuid"
)

var (
	UserInfoKey     = "INOTHCHAT:USER:"
	RoomInfoKey     = "INOTHCHAT:ROOM:"
	RoomListKey     = "INOTHCHAT:ROOMLIST"
	RoomUserListKey = "INOTHCHAT:ROOM:USERS:"
)

func RandomId() string {
	id, _ := uuid.NewUUID()
	return id.String()[:8]
}

func ToJson(obj interface{}) []byte {
	if byt, err := json.Marshal(obj); err == nil {
		return byt
	}
	return nil
}
