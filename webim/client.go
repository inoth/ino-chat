package webim

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 允许向对等方写入消息的时间。
	writeWait = 10 * time.Second

	// 允许从对等方读取下一个pong消息的时间.
	pongWait = 60 * time.Second

	// 将ping发送给具有此期间的对等方。必须小于pongWait。
	pingPeriod = (pongWait * 9) / 10

	// 对等端允许的最大邮件大小。
	maxMessageSize = 512
)

var (
	newline  = []byte{'\n'}
	space    = []byte{' '}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type User struct {
	Uid    string
	Name   string
	Avatar string
}

type Client struct {
	user *User
	cm   *ClientManager
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// 枢纽关闭了管道。
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 将排队的聊天消息添加到当前websocket消息中。
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func Register(user *User, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &Client{
		user: user,
		cm:   Instance(),
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.cm.register <- client
	// 只需要写，读取另作
	go client.write()
}
