package webim

import "github.com/gorilla/websocket"

type Client struct {
	user int
	hub  *ClientHub
	conn *websocket.Conn
	send chan []byte
}
