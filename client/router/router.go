package router

import (
	"inochat/client/config"
	col "inochat/client/controller"
	mid "inochat/client/middleware"

	"github.com/gin-gonic/gin"
)

func ServeStart() {
	r := gin.New()

	r.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})

	api := r.Group("/api")
	{
		api.POST("/init", col.InitUser)
	}
	room := api.Group("/room", mid.AuthMiddleware)
	{
		room.GET("", col.RoomList) // 房间列表
		room.POST("/msg")          // 发送消息
		room.POST("/create")       // 创建房间
		room.POST("/join/:rid")    // 加入放假
		room.POST("/exit")         // 退出房间
	}
	r.Run(config.Instance().ServerPort)
}
