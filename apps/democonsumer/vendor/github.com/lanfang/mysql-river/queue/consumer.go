package queue

import (
	"fmt"
	"github.com/bsm/sarama-cluster"
	"github.com/lanfang/go-lib/log"
	"github.com/lanfang/mysql-river/protocol"
)

type ConsumerCfg struct {
	TopicList []string
	GrouponId string
	Addr      []string
}

type Worker func(msg *protocol.EventData) error
type Consumer interface {
	Run() error
}

var consumer Consumer

func SetConsumer(x Consumer) {
	consumer = x
}
func StartConsume() error {
	if consumer == nil {
		err := fmt.Errorf("queue consumer is emtpy, need add consumer")
		log.Error("Consume %+v", err)
		return err
	}
	go consumer.Run()
	return nil
}

type SimpleConsumer struct {
	Cfg  *ConsumerCfg
	Work Worker
}

func (consumer *SimpleConsumer) Run() error {
	log.Info("Begin SimpleConsumer Run")
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.ClientID = consumer.Cfg.GrouponId
	c, err := cluster.NewConsumer(consumer.Cfg.Addr, consumer.Cfg.GrouponId, consumer.Cfg.TopicList, config)
	if err != nil {
		log.Error("Consumer NewConsumer Failed cfg:%+v, err:%+v", *consumer.Cfg, err)
		return err
	}
	defer c.Close()
	for {
		select {
		case err, ok := <-c.Errors():
			if !ok {
				log.Error("SimpleConsumer Run Errors Channel Closed cfg:%+v", *consumer.Cfg)
			} else {
				log.Error("SimpleConsumer Run Occur Err: cfg:%+v, err:%+v", *consumer.Cfg, err)
			}
		case note, ok := <-c.Notifications():
			if !ok {
				log.Info("SimpleConsumer Run Notifications Channel Closed cfg:%+v", *consumer.Cfg)
			} else {
				log.Info("SimpleConsumer Run Occur Notifications: cfg:%+v, note:%+v", *consumer.Cfg, note)
			}
		case msg, ok := <-c.Messages():
			if !ok {
				log.Error("Consumer Run Messages Channel Closed cfg:%+v", *consumer.Cfg)
				return fmt.Errorf("Messages Channel Closed")
			} else {
				data := &protocol.EventData{}
				if _, err := data.UnmarshalMsg([]byte(msg.Value)); err != nil {
					log.Error("SimpleConsumer Run Unmarshal EventData err cfg:%+v, err:%+v", *consumer.Cfg, err)
					break
				}
				if err := consumer.Work(data); err != nil {
					log.Error("SimpleConsumer Run Handle msg error cfg:%+v, data:%+v, err%+v", *consumer.Cfg, *data, err)
					//break
				}
				c.MarkOffset(msg, "")
				log.Info("SimpleConsumer msg.Topic:%+v, msg.Partition:%+v, msg.Offset:%+v", msg.Topic, msg.Partition, msg.Offset)
			}
		}
	}
	log.Info("End SimpleConsumer Run")
	return nil
}
