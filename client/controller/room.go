package controller

import (
	"inochat/client/cache"
	"inochat/client/config"
	"inochat/client/db"
	"inochat/client/db/entity"
	"inochat/client/model"
	"inochat/client/model/reply"
	"inochat/client/model/request"
	"inochat/client/res"
	"inochat/client/util"
	"inochat/client/util/nsqmsg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func RoomList(c *gin.Context) { // 换成 mongodb 中查询
	rooms, err := db.FindAll(bson.D{}, &entity.RoomInfo{})
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
	rinfo := &reply.RoomInfo{
		Rid:   rid,
		RName: req.RoomName,
	}
	uid := c.Request.Header.Get("USER_ID")

	if !db.Create(&entity.RoomInfo{
		Rid:   rid,
		RName: req.RoomName,
		Owner: uid,
	}) {
		c.JSON(500, res.Err("create room err"))
		return
	}

	if !joinRoom(uid, rid) {
		logrus.Warn("join room err")
		c.JSON(500, res.Err("join room err"))
		return
	}
	c.JSON(200, res.OK("ok", rinfo))
}

func JoinRoom(c *gin.Context) {
	var req request.RoomOperateReq
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Errorf("%v", err)
		c.JSON(400, res.ParamErr("param err"))
		return
	}

	uid := c.Request.Header.Get("USER_ID")
	// 检查房间合法性
	if db.Count(bson.D{{"rid", req.Rid}}, &entity.RoomInfo{}) == 0 {
		logrus.Warn("invalid room")
		c.JSON(400, res.ParamErr("invalid room"))
		return
	}
	if !joinRoom(uid, req.Rid) {
		logrus.Warn("join room err")
		c.JSON(500, res.Err("join room err"))
		return
	}
	c.JSON(200, res.ResultOK())
}

func ExitRoom(c *gin.Context) {
	// 退出房间
	var req request.RoomOperateReq
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Errorf("%v", err)
		c.JSON(400, res.ParamErr("param err"))
		return
	}
	uid := c.Request.Header.Get("USER_ID")
	if db.Count(bson.D{{"rid", req.Rid}}, &entity.RoomInfo{}) == 0 {
		logrus.Warn("invalid room")
		c.JSON(400, res.ParamErr("invalid room"))
		return
	}
	// 该用户退出房间
	_ = cache.SRem(config.ROOMMEMBERS+req.Rid, uid)
	// 如果当前房间内人员为 0，直接删除该房间
	if !cache.Exists(config.ROOMMEMBERS + req.Rid) {
		db.Delete(bson.D{{"rid", req.Rid}}, &entity.RoomInfo{})
		logrus.Info("房间内人数为0，默认解散")
		c.JSON(200, res.OK("房间内人数为0"))
		return
	}
	users, err := cache.SMembers(config.ROOMMEMBERS + req.Rid)
	if err != nil {
		logrus.Warn("invalid room members")
		c.JSON(400, res.ParamErr("invalid room members"))
		return
	}
	sendMsg(&model.MessageNsqBody{
		MsgType:  0, // system msg
		Target:   users,
		FromUser: uid,
		Body:     "用户" + uid + "退出房间",
	})
	c.JSON(200, res.ResultOK())
}

func RemoveRoom(c *gin.Context) {
	var req request.RoomOperateReq
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Errorf("%v", err)
		c.JSON(400, res.ParamErr("param err"))
		return
	}
	if !removeRoom(req.Rid) {
		c.JSON(500, res.Err("invalid room"))
		return
	}
	c.JSON(200, res.ResultOK())
}

func SendMsg(c *gin.Context) {
	// 给当前用户所在房间所有用户发送消息
	var req request.RoomOperateReq
	err := c.BindJSON(&req)
	if err != nil {
		logrus.Errorf("%v", err)
		c.JSON(400, res.ParamErr("param err"))
		return
	}
	uid := c.Request.Header.Get("USER_ID")
	users, err := cache.SMembers(config.ROOMMEMBERS + req.Rid)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	sendMsg(&model.MessageNsqBody{
		MsgType:  1, // user msg
		Target:   users,
		FromUser: uid,
		Body:     req.Msg,
	})
	c.JSON(200, res.ResultOK())
}

func joinRoom(uid, rid string) bool {
	err := cache.SAdd(config.ROOMMEMBERS+rid, uid)
	if err != nil {
		logrus.Errorf("%v", err)
		return false
	}
	users, err := cache.SMembers(config.ROOMMEMBERS + rid)
	if err != nil {
		logrus.Error(err.Error())
		return false
	}
	// 加入后发送一条加入信息
	sendMsg(&model.MessageNsqBody{
		MsgType:  0, // system msg
		Target:   users,
		FromUser: uid,
		Body:     "用户" + uid + "加入房间",
	})
	return true
}

func removeRoom(rid string) bool {
	if !db.Delete(bson.D{{"rid", rid}}, &entity.RoomInfo{}) {
		return false
	}
	cache.Del(config.ROOMMEMBERS + rid)
	return true
}

func sendMsg(msgBody *model.MessageNsqBody) {
	nsqmsg.CH_msg <- util.ToJson(msgBody)
}
