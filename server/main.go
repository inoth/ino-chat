package main

import (
	"inochat/server/config"
	svc "inochat/server/service"
	"inochat/server/webim"

	"github.com/gin-gonic/gin"
)

const (
	c_Topic   = "SendMsg"
	c_Channel = "WSServer"
)

func main() {
	go webim.Instance().Run()
	go svc.InitConsumer(&svc.MessageNsq{
		Topic:   c_Topic,
		Channel: c_Channel,
		Address: config.Instance().Nsq,
	})
	r := gin.New()
	r.GET("/ws", func(c *gin.Context) {
		uid := c.Query("u")
		webim.Register(uid, c.Writer, c.Request)
	})
	r.Run(config.Instance().ServerPort)
}
