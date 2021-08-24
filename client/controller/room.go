package controller

import (
	"inochat/client/cache"
	"inochat/client/config"
	"inochat/client/res"

	"github.com/gin-gonic/gin"
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

}

func JoinRoom(c *gin.Context) {

}

func ExitRoom(c *gin.Context) {
	// 退出房间，如果当前房间内人员为 0，直接删除该房间
}

func RemoveRoom(C *gin.Context) {

}

func removeRoom(rid string) {

}

func SendMsg(c *gin.Context) {
	// 给当前用户所在房间所有用户发送消息
}
