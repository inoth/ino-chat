package util

import (
	"encoding/json"

	"github.com/google/uuid"
)

func RandomID() string {
	id, _ := uuid.NewUUID()
	return id.String()
}

func ToJson(obj interface{}) []byte {
	if byt, err := json.Marshal(obj); err == nil {
		return byt
	}
	return nil
}
func ToJsonStr(obj interface{}) string {
	if byt, err := json.Marshal(obj); err == nil {
		return string(byt)
	}
	return ""
}
