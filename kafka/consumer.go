package kafka

import (
	"context"
	"demo-to-start/common"
	"encoding/base64"
	"errors"
	"github.com/Shopify/sarama"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ConsumerConfig 是 kafka consumer 相关配置
type ConsumerConfig struct {
	Brokers           []string // 必须; kafka brokers
	Version           string   // 可选; kafka 版本
	ClientID          string   // 可选; 客户端标识
	Group             string   // 必须; Consumer Group
	FromOldest        bool     // 可选; 是否从最老的记录开始读取, 默认 false
	User              string   // 可选; kafka 用户名
	Password          string   // 可选; kafka 密码
	ChannelBufferSize int      // 可选; partition consumer 缓存大小
}

// NewKafkaConsumer 创建一个新的 kafka Consumer.
//
// NOTE: 不要忘记调用 Consumer.Close, 否则会有资源泄漏.
func NewKafkaConsumer(config ConsumerConfig) (Consumer, error) {
	// 检查参数
	if len(config.Brokers) == 0 {
		return nil, errors.New("empty brokers")
	}

	// parse version
	kafkaVersion := defaultKafkaVersion
	if version := config.Version; version != "" {
		v, err := sarama.ParseKafkaVersion(version)
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

	if config.Group == "" {
		return nil, errors.New("empty group")
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

		kafkaConfig.Consumer.MaxWaitTime = time.Millisecond * 500
		kafkaConfig.Consumer.Return.Errors = true
		kafkaConfig.Consumer.Offsets.AutoCommit.Interval = time.Second
		if config.FromOldest {
			kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
		}
		if config.ChannelBufferSize > 0 {
			kafkaConfig.ChannelBufferSize = config.ChannelBufferSize
		}

		// 校验参数是否配置正确
		if err := kafkaConfig.Validate(); err != nil {
			return nil, err
		}
	}

	consumerGroup, err := sarama.NewConsumerGroup(config.Brokers, config.Group, kafkaConfig)
	if err != nil {
		return nil, err
	}

	consumer := &kafkaConsumer{
		consumerGroup: consumerGroup,
		closing:       make(chan struct{}),
	}

	// track errors
	consumer.wg.Add(1)
	go func(consumer *kafkaConsumer) {
		defer consumer.wg.Done()
		for err := range consumer.consumerGroup.Errors() {
			var (
				consumerError    sarama.ConsumerError
				ptrConsumerError *sarama.ConsumerError
				insideError      error
			)
			if errors.As(err, &ptrConsumerError) {
				insideError = ptrConsumerError.Err
			} else if errors.As(err, &consumerError) {
				insideError = consumerError.Err
			}
			if insideError != nil && errors.Is(insideError, sarama.ErrRequestTimedOut) {
				log.Println(context.Background(), "got-kafka-consume-error", "error", err.Error())
				continue
			}
			log.Println(context.Background(), "got-kafka-consume-error", "error", err.Error())
		}
	}(consumer)

	return consumer, nil
}

var (
	_ error = sarama.ConsumerError{}
	_ error = (*sarama.ConsumerError)(nil)
)

type kafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup

	started common.Bool    // 是否已经启动
	closed  common.Bool    // 是否已经关闭
	closing chan struct{}  // 关闭信号
	wg      sync.WaitGroup // 关闭之后的 sync.WaitGroup
}

func (impl *kafkaConsumer) Close(ctx context.Context) error {
	if !impl.closed.CompareAndSwap(false, true) {
		return errors.New("the consumer close method has been called")
	}
	close(impl.closing)
	if err := impl.consumerGroup.Close(); err != nil {
		log.Println(ctx, "kafka-consumer-close-failed", "error", err.Error())
		impl.wg.Wait()
		return err
	}
	impl.wg.Wait()
	log.Println(ctx, "kafka-consumer-closed")
	return nil
}

func (impl *kafkaConsumer) StartConsumeMessage(ctx context.Context, handlers map[MessageType]MessageHandler) error {
	if !impl.started.CompareAndSwap(false, true) {
		return errors.New("the consumer has been started")
	}

	if len(handlers) == 0 {
		return errors.New("empty handlers")
	}

	var groupHandler sarama.ConsumerGroupHandler = &consumerGroupHandler{
		handlers: handlers,
	}

	// 确定需要消费的 topics
	topics := make([]string, 0, len(handlers))
	for msgType := range handlers {
		topics = append(topics, kafkaTopicFromMsgType(msgType))
	}

	for {
		select {
		case <-impl.closing:
			return nil
		default:
		}

		err := impl.consumerGroup.Consume(ctx, topics, groupHandler)
		if err != nil {
			log.Println(ctx, "kafka-consume-failed", "topics", topics, "error", err.Error())
			continue
		}
	}
}

/**************************************** implements sarama.ConsumerGroupHandler ****************************************/

var _ sarama.ConsumerGroupHandler = (*consumerGroupHandler)(nil)

type consumerGroupHandler struct {
	handlers map[MessageType]MessageHandler
}

func (impl *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (impl *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (impl *consumerGroupHandler) ConsumeClaim(ss sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		err := impl.handleMessage(msg)
		_ = err // TODO: 向前消费, 不考虑重试的事情, 以后再优化这里吧
		ss.MarkMessage(msg, "")
	}
	return nil
}

func msgTypeFromKafkaTopic(topic string) (MessageType, bool) {
	if !strings.HasPrefix(topic, kafkaTopicPrefix) {
		return 0, false
	}
	n, err := strconv.ParseInt(topic[len(kafkaTopicPrefix):], 10, 32)
	if err != nil {
		return 0, false
	}
	return MessageType(n), true
}

func (impl *consumerGroupHandler) handleMessage(msg *sarama.ConsumerMessage) error {
	msgType, ok := msgTypeFromKafkaTopic(msg.Topic)
	if !ok {
		log.Println("unexpected-topic", "msg-value", string(msg.Value))
		return nil // 忽略消息, 正常情况下不会出现
	}

	// 查找 handler
	handler, ok := impl.handlers[msgType]
	if !ok || handler == nil {
		log.Println("not-found-handler", "msg-value", string(msg.Value))
		return nil // 忽略消息, 正常情况下不会出现
	}

	// base64 解码
	msgValue := make([]byte, base64.StdEncoding.DecodedLen(len(msg.Value)))
	n, err := base64.StdEncoding.Decode(msgValue, msg.Value)
	if err != nil {
		log.Println("base64-decode-msg-failed", "msg-value", string(msg.Value), "error", err.Error())
		return nil // 忽略消息, 正常情况下不会出现
	}
	msgValue = msgValue[:n]

	// 处理消息
	bizMsg := &Message{
		Value: msgValue,
		Kafka: MessageForKafka{
			Timestamp:      msg.Timestamp,
			BlockTimestamp: msg.BlockTimestamp,
			Topic:          msg.Topic,
			Partition:      msg.Partition,
			Offset:         msg.Offset,
		},
	}
	err = handler.ServeMessage(context.Background(), bizMsg)
	if err != nil {
		log.Println("handle-kafka-message-bus-message-failed", "msg-value", string(msg.Value), "error", err.Error())
		return err
	}
	return nil
}
