package util

import "github.com/google/uuid"

func RandomID() string {
	id, _ := uuid.NewUUID()
	return id.String()
}
