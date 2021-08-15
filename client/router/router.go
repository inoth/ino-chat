package router

import "github.com/gin-gonic/gin"

func InitClient() {
	r := gin.New()
	api := r.Group("/api")
	{
		api.POST("/init", func(c *gin.Context) {
			c.String(200, "init user")
		})
		api.POST("/room/create")
		api.POST("/room/join/:rid")
		api.POST("/room/exit")
		api.POST("/msg/send")
	}
	r.Run(":9978")
}
