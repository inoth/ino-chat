package controller

import (
	"ino-chat/res"
	"ino-chat/src/models"
	svc "ino-chat/src/service"
	"ino-chat/util"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// @Summary 创建房间
// @Tags 房间相关
// @Accept  json
// @Produce  json
// @Param object body models.NewRoomBody true "创建房间body"
// @Success 200 object res.ApiResult
// @Router /api/room/new [post]
func NewRoom(c *gin.Context) {
	var req models.NewRoomBody
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, res.ParamErr("param err."))
		return
	}
	rid := util.RandomId()
	uid := c.GetHeader("USER_ID")
	if !svc.NewRoom(rid, uid, req.Title) {
		c.JSON(500, res.Err("room create err."))
		return
	}
	c.JSON(200, res.OK("ok", &models.RoomInfo{
		Rid:   rid,
		Owner: uid,
		Title: req.Title,
	}))
}

// @Summary 查询房间列表
// @Tags 房间相关
// @Produce  json
// @Success 200 object res.ApiResult
// @Router /api/room [get]
func GetRooms(c *gin.Context) {
	rooms, err := svc.GetAllRoom()
	if err != nil {
		log.Error(err.Error())
		c.JSON(404, res.NotFound("not found."))
		return
	}
	c.JSON(200, res.OK("ok", rooms))
}

// @Summary 加入房间
// @Tags 房间相关
// @Produce  json
// @Param Authorization header string false "用户令牌"
// @Param rid path string true "房间id"
// @Success 200 object res.ApiResult
// @Router /api/room/{rid}/join [post]
func JoinRoom(c *gin.Context) {
	rid := c.Param("rid")
	uid := c.GetHeader("USER_ID")
	if !svc.JoinRoom(rid, uid) {
		c.JSON(500, res.Err("room join err."))
		return
	}
	c.JSON(200, res.ResultOK())
}

// @Summary 退出房间
// @Tags 房间相关
// @Produce  json
// @Param Authorization header string false "用户令牌"
// @Param rid path string true "房间id"
// @Success 200 object res.ApiResult
// @Router /api/room/exit/{rid} [post]
func ExitRoom(c *gin.Context) {
	rid := c.Param("rid")
	uid := c.GetHeader("USER_ID")
	if !svc.ExitRoom(rid, uid) {
		c.JSON(500, res.Err("room join err."))
		return
	}
	svc.MessageRoomPush(uid, rid, "system", "用户："+uid+"退出了房间")
	c.JSON(200, res.ResultOK())
}
