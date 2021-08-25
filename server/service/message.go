package service

import (
	"encoding/json"
	"inochat/server/cache"
	"inochat/server/config"
	"inochat/server/webim"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type RegisteredConsumer interface {
	GetTopic() string
	GetChannel() string
	GetAddress() string
}
type MessageNsq struct {
	Topic,
	Channel,
	Address string
}

func (c *MessageNsq) GetTopic() string   { return c.Topic }
func (c *MessageNsq) GetChannel() string { return c.Channel }
func (c *MessageNsq) GetAddress() string { return c.Address }

type MessageNsqBody struct {
	MsgType    int    `json:"msgType"`
	TargetType int    `json:"targetType"`
	Target     string `json:"target"`
	FromUser   string `json:"fromUser"`
	Body       string `json:"body"`
}

func (c *MessageNsq) HandleMessage(msg *nsq.Message) error {
	data := &MessageNsqBody{}
	if err := json.Unmarshal(msg.Body, data); err != nil {
		return err
	}
	target := make([]string, 0)
	switch data.TargetType {
	case 0: // 单对单发送
		target = append(target, data.Target)
	case 1: // 发送至房间
		users, err := cache.SMembers(config.ROOMMEMBERS + data.Target)
		if err != nil {
			logrus.Error(err.Error())
			break
		}
		target = append(target, users...)
	}
	if len(target) <= 0 {
		return nil
	}
	err := webim.SendMessage(int32(data.MsgType), time.Now().Unix(), target, data.FromUser, data.Body)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func InitConsumer(config RegisteredConsumer) {
	cfg := nsq.NewConfig()
	c, err := nsq.NewConsumer(config.GetTopic(), config.GetChannel(), cfg)
	if err != nil {
		panic(err)
	}
	c.AddHandler(config.(nsq.Handler))

	if err := c.ConnectToNSQD(config.GetAddress()); err != nil {
		panic(err)
	}
}
