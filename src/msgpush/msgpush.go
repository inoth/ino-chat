package msgpush

import (
	"ino-chat/config"

	"github.com/nsqio/go-nsq"
	log "github.com/sirupsen/logrus"
)

const (
	SendTopic = ""
)

var ch_msg = make(chan []byte, 10)

func SendMessage(msg []byte) {
	ch_msg <- msg
}

func StarMsgPushSvc() {
	pro, err := nsq.NewProducer(config.Instance().Nsq, nsq.NewConfig())
	if err != nil {
		log.Errorf("%v", err)
		log.Panic(err.Error())
	}
	log.Info("connect nsq.")
	defer func() {
		defer pro.Stop()
		close(ch_msg)
		if err := recover(); err != nil {
			log.Errorf("%v", err)
		}
	}()
	// 持续捕获发送的消息
	for {
		select {
		case msg := <-ch_msg:
			err := pro.Publish(SendTopic, msg)
			if err != nil {
				log.Warn("msg send err")
			}
		}
	}
}
