package chat

import (
	"ino-chat/src/models"
	"sync"
)

var (
	hub  *Hub
	once sync.Once
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	// broadcast chan []byte
	broadcast chan *models.MsgPushBody

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan string
}

func NewHub() *Hub {
	once.Do(func() {
		hub = &Hub{
			broadcast:  make(chan *models.MsgPushBody),
			register:   make(chan *Client),
			unregister: make(chan string),
			clients:    make(map[string]*Client),
		}
	})
	return hub
}

func SendMsg(msg *models.MsgPushBody) {
	hub.broadcast <- msg
}

func Disconnect(uid string) {
	hub.unregister <- uid
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.uid] = client

		case uid := <-h.unregister:
			if client, ok := h.clients[uid]; ok {
				delete(h.clients, uid)
				close(client.send)
			}

		case message := <-h.broadcast:
			// 仅给连接当前节点用户发送消息
			for _, target := range message.Target {
				if client, ok := h.clients[target]; ok {
					data := &models.WsSendMessage{
						FromUid: message.FromUid,
						MsgType: message.MsgType,
						Body:    message.Body,
					}
					client.send <- data.ToJson()
				}
			}
			// for uid, client := range h.clients {
			// 	select {
			// 	case client.send <- []byte(data.Body):
			// 	default:
			// 		close(client.send)
			// 		delete(h.clients, uid)
			// 	}
			// }
		}
	}
}
