package controller

import (
	"ino-chat/res"
	"ino-chat/src/models"
	"ino-chat/src/service"
	svc "ino-chat/src/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SendMessage(c *gin.Context) {
	var req models.MsgSendBody
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, res.ParamErr("param err."))
		return
	}
	uid := c.GetHeader("USER_ID")
	switch req.TargetType {
	case "room":
		service.MessageRoomPush(uid, req.Target, req.MsgType, req.Msg)
	case "user":
		service.MessagePush(uid, req.Target, req.MsgType, req.Msg)
	default:
		logrus.Warnf("target [%v] err.", req.TargetType)
	}
	// c.JSON(200, res.ResultOK())
}

func LoadChatPage(c *gin.Context) {
	uid := c.GetHeader("USER_ID")
	// default join testroom
	svc.JoinRoom("testroom", uid)
	c.HTML(200, "index.html", "ws://"+c.Request.Host+"/ws")
}
