package mq_consumer

import (
	"context"
	"core/biz/socket"
	"core/conf"
	chat "core/hertz_gen/chat"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

var c = conf.GetConf()

func StartConsumer() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   c.Kafka.Address,
		Topic:     c.Kafka.Topic,
		GroupID:   "chat-consumer",
		Partition: 0, // 可选，用 GroupID 自动分配分区更常用
		MinBytes:  1,
		MaxBytes:  10e6,
	})
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			hlog.Error("Kafka read error:", err)
			continue
		}
		var msg chat.ChatMsg
		if err := proto.Unmarshal(m.Value, &msg); err != nil {
			hlog.Error("Unmarshal error:", err)
			continue
		}
		if err := socket.SaveTextMsg(&msg); err != nil {
			hlog.Error("SaveTextMsg error:", err)
		}
	}
}
