package service

import (
	"encoding/json"

	"github.com/nsqio/go-nsq"
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
	MsgType    int
	TargetType int
	Target     string
	FromUser   string
	Body       string
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

	}
	// webim.SendMessage()
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
