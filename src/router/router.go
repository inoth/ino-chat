package router

import (
	"ino-chat/config"
	ex "ino-chat/exception"
	"ino-chat/src/chat"
	col "ino-chat/src/controller"
	mid "ino-chat/src/middleware"
	"net/http"

	"ino-chat/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ServerStar() {
	r := gin.New()
	r.Use(ex.ExceptionHandle)

	r.MaxMultipartMemory = 10 << 20

	isDev := gin.Mode()
	if isDev != "release" {
		docs.SwaggerInfo.BasePath = ""
		ginSwagger.WrapHandler(swaggerfiles.Handler,
			ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
			ginSwagger.DefaultModelsExpandDepth(-1))

		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	api := r.Group("/api")
	{
		api.POST("login", col.Login)
		// user := api.Group("/user", mid.AuthMiddleware)
		// {
		// 	user.GET("")
		// }
		room := api.Group("/room", mid.AuthMiddleware)
		{
			room.GET("", col.GetRooms)
			room.POST("/new", col.NewRoom)
			room.POST("/join/:rid", col.JoinRoom)
			room.POST("/exit/:rid", col.ExitRoom)
		}
		msg := api.Group("/msg", mid.AuthMiddleware)
		{
			msg.POST("/send", col.SendMessage)
		}
	}

	// chat web page
	{
		r.LoadHTMLGlob("template/*")
		r.GET("/", func(c *gin.Context) {
			c.HTML(200, "login.html", "")
		})
		r.GET("/chat", mid.AuthMiddleware, col.LoadChatPage)
	}

	// ws server start
	{
		hub := chat.NewHub()
		go hub.Run()
		r.GET("/ws", mid.AuthMiddleware, func(c *gin.Context) {
			uid := c.GetHeader("USER_ID")
			logrus.Infof("%v link", uid)
			chat.ServeWs(uid, hub, c.Writer, c.Request)
		})
	}

	r.Run(config.Instance().ServerPort)
}
