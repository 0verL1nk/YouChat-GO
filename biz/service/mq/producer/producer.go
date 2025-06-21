package mq_producer

import (
	"context"
	"core/conf"
	"core/hertz_gen/chat"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"google.golang.org/protobuf/proto"

	"github.com/segmentio/kafka-go"
)

var c = conf.GetConf()
var writer = kafka.NewWriter(
	kafka.WriterConfig{
		Brokers:  c.Kafka.Address,
		Topic:    c.Kafka.Topic,
		Balancer: &kafka.Hash{},
	})

func HandlerWSMessage(msg *chat.ChatMsg) {
	data, _ := proto.Marshal(msg)
	key := msg.To
	err := writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key),
		Value: data,
	})
	if err != nil {
		hlog.Error("Kafka write error:", err)
	}
}
