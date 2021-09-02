package nsqmsg

import (
	"inochat/client/config"

	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

const (
	SENDTOPIC = "SendMsg"
)

var CH_msg = make(chan []byte, 10)

func InitProducer() {
	pro, err := nsq.NewProducer(config.Instance().Nsq, nsq.NewConfig())
	if err != nil {
		logrus.Errorf("%v", err)
		logrus.Panic(err.Error())
	}
	logrus.Info("connect nsq.")
	defer func() {
		defer pro.Stop()
		close(CH_msg)
		if err := recover(); err != nil {
			logrus.Errorf("%v", err)
		}
	}()
	// 持续捕获发送的消息
	for {
		select {
		case msg := <-CH_msg:
			err := pro.Publish(SENDTOPIC, msg)
			if err != nil {
				logrus.Warn("msg send err")
			}
		default:
		}
	}
}
