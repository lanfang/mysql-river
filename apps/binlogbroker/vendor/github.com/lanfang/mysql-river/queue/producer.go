package queue

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/lanfang/go-lib/log"
	"github.com/lanfang/mysql-river/protocol"
	syslog "log"
	"os"
	"time"
)

func init() {
	sarama.Logger = syslog.New(os.Stdout, "[Sarama] ", syslog.LstdFlags|syslog.Lshortfile)
}

var defaultQueue KafKaQueue

type KafKaQueue struct {
	Addr     []string
	Producer sarama.SyncProducer
}

func Init(addr []string) error {
	log.Info("Begin Queue Init addr:%+v", addr)
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 2
	config.Version = sarama.V0_10_0_0
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Flush.MaxMessages = 1
	config.Producer.Return.Successes = true

	defaultQueue = KafKaQueue{Addr: addr}
	var err error
	if defaultQueue.Producer, err = sarama.NewSyncProducer(defaultQueue.Addr, config); err != nil {
		log.Error("Queue Init Failed addr:%+v, err:%+v", addr, err)
	}
	log.Info("Queue Init End addr:%+v, err:%+v", addr, err)
	return err
}

func PushBack(data *protocol.EventData) error {
	msg := genKafkaMsg(data)
	var err = retrySendMessages(data, msg)
	if err != nil {
		log.Error("Queue PushBack err:%v, data:%+v", err, *data)
	} else {
		log.Info("push to kafka [%v %v.%v] success", data.Action, data.Schema, data.Table)
	}
	return err
}

func genKafkaMsg(data *protocol.EventData) []*sarama.ProducerMessage {
	msg_list := make([]*sarama.ProducerMessage, 0)
	if data.Owner.KeyIndex >= 0 {
		dispatch_msg := make(map[protocol.TopicInfo]*protocol.EventData)
		for _, row := range data.Rows {
			data.Owner.Key = fmt.Sprintf("%v", row[data.Owner.KeyIndex])
			if _, ok := dispatch_msg[data.Owner]; !ok {
				dispatch_msg[data.Owner] = &protocol.EventData{
					Action: data.Action, Schema: data.Schema, Table: data.Table,
					Columns: data.Columns, Rows: make([][]interface{}, 0),
				}
			}
			dispatch_msg[data.Owner].Rows = append(dispatch_msg[data.Owner].Rows, row)
		}
		for topic, data := range dispatch_msg {
			msg := &sarama.ProducerMessage{
				Topic: topic.Topic, Key: sarama.StringEncoder(topic.Key), Value: data}
			msg_list = append(msg_list, msg)
		}
	} else {
		msg := &sarama.ProducerMessage{
			Topic: data.Owner.Topic, Key: sarama.StringEncoder(data.Owner.Key), Value: data}
		msg_list = append(msg_list, msg)
	}
	return msg_list
}

func retrySendMessages(data *protocol.EventData, msgs []*sarama.ProducerMessage) error {
	var err error
	for loop := true; loop; loop = false {
		if err = defaultQueue.Producer.SendMessages(msgs); err == nil {
			break
		}
		log.Info("begin retry push msg, last err:%+v ,data:%+v", err, *data)
		deadline := time.After(2 * time.Minute)
		interval := 5 * time.Second
		tick := time.NewTicker(interval)
		defer tick.Stop()
	InnerLoop:
		for {
			select {
			case <-tick.C:
				log.Info("retry push msg, last err:%+v ,data:%+v", err, *data)
				if err = defaultQueue.Producer.SendMessages(msgs); err == nil {
					break InnerLoop
				}
			case <-deadline:
				log.Info("retry timeout with err:%v, data:%+v", err, *data)
				break InnerLoop
			}
		}
	}
	return err
}
