package router

import "github.com/gin-gonic/gin"

func InitClient() {
	r := gin.New()
	api := r.Group("/api")
	{
		api.POST("/init", func(c *gin.Context) {
			c.String(200, "init user")
		})
		api.POST("/msg/send")
	}
	room := api.Group("/room")
	{
		room.POST("/create")
		room.POST("/join/:rid")
		room.POST("/exit")
	}
	r.Run(":9978")
}
