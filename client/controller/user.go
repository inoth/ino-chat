package controller

import (
	mid "inochat/client/middleware"
	"inochat/client/model/reply"
	"inochat/client/model/request"
	"inochat/client/res"
	"inochat/client/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func InitUser(c *gin.Context) {
	var req request.UserInitReq
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Errorf("%v", err)
		c.JSON(res.PARAMETERERR, res.ParamErr("param err"))
		return
	}
	uid := util.RandomID()
	sign, err := mid.CreateToken(uid, req.UserName)
	if err != nil {
		logrus.Errorf("%v", err)
		c.JSON(500, res.Err(err.Error()))
	}
	c.Header("Authorization", sign)
	c.SetCookie("Authorization", sign, 60*60*24, "/", "api.inoth.site", false, true)
	c.JSON(200, res.OK("ok", &reply.UserInfo{
		Uid:      uid,
		UserName: req.UserName,
	}))
}
