package main

import (
	"inochat/client/router"
	"inochat/client/util/nsqmsg"
)

func main() {
	go nsqmsg.InitProducer()
	router.ServeStart()
}
