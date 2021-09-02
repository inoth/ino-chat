package entity

type IEntity interface {
	Col() string
}
type RoomInfo struct {
	Rid   string `json:"rid"`
	RName string `json:"rname"`
	Owner string `json:"owner"`
}

func (m *RoomInfo) Col() string {
	return "room_info"
}
