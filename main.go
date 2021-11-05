package main

import (
	"ino-chat/cache"
	"ino-chat/src/router"
)

func main() {
	cache.Init()
	// go mp.StarMsgPushSvc()
	router.ServerStar()
}
