package request

type UserInitReq struct {
	UserName string `json:"userName"`
}

type RoomCreateReq struct {
	RoomName string `json:"roomName"`
}
