package kafka

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
	"unsafe"
)

// defaultKafkaVersion 默认的 kafka 版本号, 和阿里云 kafka 版本保持一致.
var defaultKafkaVersion = sarama.V2_2_0_0
var (
	timeLayout = "2006-01-02 15:04:05"
)

type Producer interface {
	// SendMessage 发送一个消息到消息总线.
	//
	// ⚠️注意: message.MessageType 和 proto.Message 要匹配.
	SendMessage(context.Context, MessageType, proto.Message, ...SendMessageOption) error

	// Close 关闭 Producer, 释放相关资源, 防止资源泄漏.
	Close(context.Context) error
}

func ToJsonString(data interface{}) string {
	bs, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return *(*string)(unsafe.Pointer(&bs))
}
