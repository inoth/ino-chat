package controller

import (
	"inochat/client/cache"
	"inochat/client/config"
	"inochat/client/model"
	"inochat/client/model/request"
	"inochat/client/res"
	"inochat/client/util"
	"inochat/client/util/nsqmsg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RoomList(c *gin.Context) {
	rooms, err := cache.SMembers(config.ROOMS)
	if err != nil {
		c.JSON(404, res.NotFound("not found"))
		return
	}
	c.JSON(200, res.OK("ok", rooms))
}

func CreateRoom(c *gin.Context) {
	var req request.RoomCreateReq
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Errorf("%v", err)
		c.JSON(400, res.ParamErr("param err"))
		return
	}
	rid := util.RandomID()
	rinfo := &model.RoomInfo{
		Rid:   rid,
		RName: req.RoomName,
	}
	uid := c.Request.Header.Get("USER_ID")
	cache.SAdd(config.ROOMS, string(util.ToJson(rinfo)))
	if !joinRoom(uid, rid) {
		logrus.Warn("join room err")
		c.JSON(500, res.Err("join room err"))
		return
	}
	c.JSON(200, res.OK("ok", rinfo))
}

func JoinRoom(c *gin.Context) {
	rid := c.PostForm("rid")
	uid := c.Request.Header.Get("USER_ID")
	if !joinRoom(uid, rid) {
		logrus.Warn("join room err")
		c.JSON(500, res.Err("join room err"))
		return
	}
	c.JSON(200, res.OK("ok"))
}

func ExitRoom(c *gin.Context) {
	// 退出房间，如果当前房间内人员为 0，直接删除该房间
}

func RemoveRoom(C *gin.Context) {

}

func SendMsg(c *gin.Context) {
	// 给当前用户所在房间所有用户发送消息
}

func joinRoom(uid, rid string) bool {
	err := cache.LPush(config.ROOMMEMBERS+rid, uid)
	if err != nil {
		logrus.Errorf("%v", err)
		return false
	}
	return true
}

func removeRoom(rid string) {

}

func sendMsg(msgBody *model.MessageNsqBody) {
	nsqmsg.CH_msg <- util.ToJson(msgBody.Body)
}
