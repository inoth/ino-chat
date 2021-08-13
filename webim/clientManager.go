package webim

import (
	"sync"
	"webchat/webim/msgtype"
)

var (
	once sync.Once
	cm   *ClientManager
)

type ClientManager struct {
	// 根据 userid 保存所有连接
	clients map[string]*Client
	// 房间内多个用户
	room map[string][]*Client

	// 新增连接用户
	register chan *Client
	// 断开连接用户
	unregister chan *Client

	broadcast chan *MsgBody
}

func Instance() *ClientManager {
	once.Do(func() {
		cm = &ClientManager{
			clients:    make(map[string]*Client),
			room:       make(map[string][]*Client),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			broadcast:  make(chan *MsgBody),
		}
	})
	return cm
}

func (cm *ClientManager) JoinRoom(rid, uid string) {
	_, ok := cm.room[rid]
	if !ok {
		cm.room[rid] = make([]*Client, 0)
	}
	cm.room[rid] = append(cm.room[rid], cm.clients[uid])
}

func (cm *ClientManager) Run() {
	for {
		select {
		case client := <-cm.register:
			cm.clients[client.user.Uid] = client
		case client := <-cm.unregister:
			delete(cm.clients, client.user.Uid)
			close(client.send)
		case message := <-cm.broadcast:
			switch message.msgType {
			case msgtype.MSG: // 用户一对一发送消息
				client := cm.clients[message.target[0]]
				client.send <- message.ToJson()
			case msgtype.ROOM: // 房间内发送消息
				for _, client := range cm.room[message.target[0]] {
					client.send <- message.ToJson()
				}
			case msgtype.MSGS: // 多个用户推送消息
				for _, uid := range message.target {
					cm.clients[uid].send <- message.ToJson()
				}
			case msgtype.ROOMS: // 多个房间推送消息
				for _, rid := range message.target {
					for _, client := range cm.room[rid] {
						client.send <- message.ToJson()
					}
				}
			}
		}
	}
}
