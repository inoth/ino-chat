package main

import (
	"fmt"
	"webchat/webim"

	"github.com/gin-gonic/gin"
)

func main() {

	ws := webim.Instance()
	go ws.Run()

	r := gin.New()
	r.GET("/ws", func(c *gin.Context) {
		uid := c.Query("uid")
		name := c.Query("name")
		webim.Register(&webim.User{Uid: uid, Name: name}, c.Writer, c.Request)
	})

	r.GET("/join/:rid/:uid", func(c *gin.Context) {
		rid := c.Param("rid")
		uid := c.Param("uid")
		webim.Instance().JoinRoom(rid, uid)
		c.String(200, "ok: "+rid+" "+uid)
	})

	r.POST("send", func(c *gin.Context) {
		uid := c.PostForm("uid")
		target := c.PostForm("target")
		msg := c.PostForm("msg")
		err := webim.SendUserMsg(uid, msg, target)
		if err != nil {
			fmt.Printf("%v", err)
			c.String(500, err.Error())
			return
		}
		c.String(200, "ok")
	})

	r.Run(":9999")
}
