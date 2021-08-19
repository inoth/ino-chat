package router

import (
	col "inochat/client/controller"

	"github.com/gin-gonic/gin"
)

func ServeStart() {
	r := gin.New()
	api := r.Group("/api")
	{
		api.POST("/init", func(c *gin.Context) {
			c.String(200, "init connect.")
		})
	}
	room := api.Group("/room")
	{
		room.GET("", col.RoomList) // 房间列表
		room.POST("/msg")          // 发送消息
		room.POST("/create")       // 创建房间
		room.POST("/join/:rid")    // 加入放假
		room.POST("/exit")         // 退出房间
	}
	r.Run(":9978")
}
