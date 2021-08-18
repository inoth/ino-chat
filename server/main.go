package main

import (
	"inochat/server/config"
	"inochat/server/webim"

	"github.com/gin-gonic/gin"
)

func main() {
	go webim.Instance().Run()
	r := gin.Default()
	r.GET("/ws", func(c *gin.Context) {
		uid := c.Query("u")
		webim.Register(uid, c.Writer, c.Request)
	})
	r.Run(config.Instance().ServerPort)
}
