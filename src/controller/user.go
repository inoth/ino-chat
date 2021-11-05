package controller

import (
	"ino-chat/res"
	"ino-chat/src/middleware"
	"ino-chat/src/service"
	"ino-chat/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 登录接口
// @Tags 登录
// @Accept  json
// @Produce  json
// @Success 200 object res.ApiResult
// @Router /api/login [post]
func Login(c *gin.Context) {
	// var req models.LoginData
	// if err := c.BindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, res.ParamErr("param err."))
	// 	return
	// }
	// 插入记录进入缓存
	uid := util.RandomId()
	if !service.SaveUserInCache(uid) {
		c.JSON(500, res.Err("login err."))
		return
	}
	sign, _ := middleware.CreateToken(uid, uid)
	c.Header("Authorization", sign)
	c.SetCookie("Authorization", sign, 60*60*24, "/", "", false, true)
	c.JSON(http.StatusOK, res.OK("ok", uid))
}
