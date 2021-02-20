package kafka

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
	"log"
	"strconv"
	"time"
)

// 10.105.11.29:9092
var (
	producerClosed = 1
)

// KafkaProducerConfig 是 kafka producer 相关配置.
type KafkaProducerConfig struct {
	Brokers           []string // 必须; kafka brokers
	Version           string   // 可选; kafka 版本
	ClientID          string   // 可选; 客户端标识
	User              string   // 可选; kafka 用户名
	Password          string   // 可选; kafka 密码
	DisableLogMessage bool     // 可选; 不打印消息日志, 默认为 false, 即表示打印
}

type kafkaProducer struct {
	logMessage bool
	producer   sarama.SyncProducer
	closed     int
}

func (impl *kafkaProducer) Close(ctx context.Context) error {
	if impl.closed == producerClosed {
		return errors.New("the producer close method has been called")
	}
	impl.closed = producerClosed
	if err := impl.producer.Close(); err != nil {
		log.Println(ctx, "failed to close kafka producer", "error", err.Error())
		return err
	}
	return nil
}

type sendMessageOptions struct {
	partitionKey string
}

const kafkaTopicPrefix = "topic_"

type MessageType int64

func kafkaTopicFromMsgType(msgType MessageType) string {
	return kafkaTopicPrefix + strconv.FormatInt(int64(msgType), 10)
}
func (m *MessageType) String() string {
	if m == nil {
		return ""
	}
	return strconv.FormatInt(int64(*m), 10)
}

// SendMessageOption 发送消息的可选配置.
type SendMessageOption func(*sendMessageOptions)

func (impl *kafkaProducer) SendMessage(ctx context.Context, msgType MessageType, msg proto.Message, opts ...SendMessageOption) error {
	if impl.closed == producerClosed {
		return errors.New("the producer has been closed")
	}

	// topic
	topic := kafkaTopicFromMsgType(msgType)

	// key
	var (
		keyString string
		key       sarama.Encoder
	)
	if len(opts) > 0 {
		var o sendMessageOptions
		for _, opt := range opts {
			if opt == nil {
				continue
			}
			opt(&o)
		}
		switch {
		case o.partitionKey != "":
			keyString = o.partitionKey
			key = sarama.StringEncoder(keyString)
		}
	}

	// value
	msgData, err := Marshal(msg)
	if err != nil {
		return err
	}
	base64MsgData := make([]byte, base64.StdEncoding.EncodedLen(len(msgData)))
	base64.StdEncoding.Encode(base64MsgData, msgData)
	value := sarama.ByteEncoder(base64MsgData)

	// publish to kafka
	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   key,
		Value: value,
	}
	partition, offset, err := impl.producer.SendMessage(kafkaMsg)
	if err != nil {
		log.Println(ctx, "failed-to-send-message-to-kafka-message-bus", "msg_type", msgType.String(), "message", ToJsonString(msg), "error", err.Error())
		return err
	}

	fields := []interface{}{"msg_type", msgType.String(), "message", ToJsonString(msg), "msg_key", keyString, "msg_topic", topic, "msg_partition", partition, "msg_offset", offset}
	if !kafkaMsg.Timestamp.IsZero() {
		fields = append(fields, "msg_timestamp", kafkaMsg.Timestamp.Format(timeLayout))
	}
	if impl.logMessage {
		log.Println("success-to-send-message-to-kafka-message-bus", fields)
	}
	return nil
}

// NewKafkaProducer 创建一个新的 kafka Producer.
//
// NOTE: 不要忘记调用 Producer.Close, 否则会有资源泄漏.
func NewKafkaProducer(config KafkaProducerConfig) (Producer, error) {
	if len(config.Brokers) == 0 {
		return nil, errors.New("empty brokers")
	}

	kafkaVersion := defaultKafkaVersion
	if config.Version != "" {
		v, err := sarama.ParseKafkaVersion(config.Version)
		if err != nil {
			return nil, errors.New("invalid version")
		}
		kafkaVersion = v
	}

	// check version
	if !kafkaVersion.IsAtLeast(sarama.V0_10_2_0) {
		return nil, errors.New("version required at least 0.10.2.0")
	}

	if config.ClientID == "" {
		config.ClientID = "golang"
	}

	kafkaConfig := sarama.NewConfig()
	{
		kafkaConfig.Version = kafkaVersion

		if config.ClientID != "" {
			kafkaConfig.ClientID = config.ClientID
		}

		if config.User != "" {
			kafkaConfig.Net.SASL.Enable = true
			kafkaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
			kafkaConfig.Net.SASL.User = config.User
			kafkaConfig.Net.SASL.Password = config.Password
		}

		kafkaConfig.Net.KeepAlive = 30 * time.Second
		kafkaConfig.Producer.RequiredAcks = sarama.WaitForLocal
		kafkaConfig.Producer.Compression = sarama.CompressionSnappy
		kafkaConfig.Producer.Return.Successes = true
		kafkaConfig.Producer.Return.Errors = true

		// 校验参数是否配置正确
		if err := kafkaConfig.Validate(); err != nil {
			return nil, err
		}
	}

	producer, err := sarama.NewSyncProducer(config.Brokers, kafkaConfig)
	if err != nil {
		return nil, err
	}
	return &kafkaProducer{
		logMessage: !config.DisableLogMessage,
		producer:   producer,
	}, nil
}
