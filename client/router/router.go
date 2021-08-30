package router

import (
	"inochat/client/config"
	col "inochat/client/controller"
	mid "inochat/client/middleware"
	ex "inochat/client/middleware/exception"

	"github.com/gin-gonic/gin"
)

func ServeStart() {
	r := gin.New()

	r.Use(ex.ExceptionHandle)

	r.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})

	api := r.Group("/api")
	{
		api.POST("/init", col.InitUser)
	}
	room := api.Group("/room", mid.AuthMiddleware)
	{
		room.GET("", col.RoomList)           // 房间列表
		room.POST("/msg", col.SendMsg)       // 发送消息
		room.POST("/create", col.CreateRoom) // 创建房间
		room.POST("/join", col.JoinRoom)     // 加入放假
		room.POST("/exit", col.ExitRoom)     // 退出房间
	}
	r.Run(config.Instance().ServerPort)
}
