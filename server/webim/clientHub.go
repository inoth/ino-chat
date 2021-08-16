package webim

import "sync"

var (
	once sync.Once
	cm   *ClientHub
)

type ClientHub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *MsgBody
}

func Instance() *ClientHub {
	once.Do(func() {
		cm = &ClientHub{
			clients:    make(map[string]*Client),
			register:   make(chan *Client, 5),
			unregister: make(chan *Client, 5),
			broadcast:  make(chan *MsgBody, 10),
		}
	})
	return cm
}

func (hub *ClientHub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client.user] = client
		case client := <-hub.unregister:
			delete(hub.clients, client.user)
			close(client.send)
		case message := <-hub.broadcast:
			sendMessage(message)
		}
	}
}

func sendMessage(msg *MsgBody) {

}
