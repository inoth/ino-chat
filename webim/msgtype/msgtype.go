package msgtype

const (
	// 一对一消息
	MSG = iota
	// 房间内消息
	ROOM
	// 多个用户消息推送
	MSGS
	// 多个房间消息推送
	ROOMS
)
