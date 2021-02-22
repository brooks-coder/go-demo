package kafka

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
	"time"
	"unsafe"
)

// defaultKafkaVersion 默认的 kafka 版本号, 和阿里云 kafka 版本保持一致.
var defaultKafkaVersion = sarama.V2_2_0_0
var (
	timeLayout = "2006-01-02 15:04:05"
)

// Producer 是消息总线的生产者接口.
type Producer interface {
	// SendMessage 发送一个消息到消息总线.
	//
	// ⚠️注意: message.MessageType 和 proto.Message 要匹配.
	SendMessage(context.Context, MessageType, proto.Message, ...SendMessageOption) error

	// Close 关闭 Producer, 释放相关资源, 防止资源泄漏.
	Close(context.Context) error
}

// Consumer 是消息总线的消费者接口.
type Consumer interface {
	// StartConsumeMessage 启动消费消息总线上的消息.
	//  StartConsumeMessage 会阻塞当前的 goroutine, 直到 Close 方法被调用了.
	StartConsumeMessage(_ context.Context, handlers map[MessageType]MessageHandler) error

	// Close 停止消费消息, 关闭 Consumer, 释放相关资源, 防止资源泄漏.
	//
	// ⚠️注意: 即使没有调用 StartConsumeMessage 也需要调用这个方法, 否则有资源泄漏.
	Close(context.Context) error
}

// MessageHandler 是消息处理接口, Consumer 的实现需要支持这个接口.
type MessageHandler interface {
	// ServeMessage 处理消息总线出来的消息, 处理成功返回 nil, 否则返回相应的错误.
	//
	//  对于 mns 消息总线如果返回错误则消息不会被删除, 这个消息还会被消费到
	//  对于 kafka 消息总线如果返回错误消息则只是打印日志, 不会被再次消费到
	ServeMessage(ctx context.Context, msg *Message) error
}

// Message 是消息数据结构
type Message struct {
	Value []byte          // 具体消息的值
	MNS   MessageForMNS   // 阿里云mns消息总线特有的值
	Kafka MessageForKafka // kafka消息总线特有的值
}

// MessageForMNS mns消息
type MessageForMNS struct {
	MessageID        string // 消息编号，在一个Queue中唯一
	EnqueueTime      int64  // 消息发送到队列的时间，从1970年1月1日0点整开始的毫秒数
	FirstDequeueTime int64  // 第一次被消费的时间，从1970年1月1日0点整开始的毫秒数
	DequeueCount     int    // 总共被消费的次数
	Priority         int    // 消息的优先级权值
}

// MessageForKafka kafka消息
type MessageForKafka struct {
	Timestamp      time.Time // only set if kafka is version 0.10+, inner message timestamp
	BlockTimestamp time.Time // only set if kafka is version 0.10+, outer (compressed) block timestamp
	Topic          string
	Partition      int32
	Offset         int64
}

// ToJsonString json字符串
func ToJsonString(data interface{}) string {
	bs, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return *(*string)(unsafe.Pointer(&bs))
}

// Marshal is same as proto.Marshal(pb)
func Marshal(pb proto.Message) ([]byte, error) {
	return proto.MarshalOptions{
		AllowPartial: true, // syntax = "proto3";
	}.Marshal(pb)
}

// Unmarshal is same as proto.Unmarshal(buf,pb)
func Unmarshal(buf []byte, pb proto.Message) error {
	return proto.UnmarshalOptions{
		AllowPartial: true, // syntax = "proto3";
	}.Unmarshal(buf, pb)
}
