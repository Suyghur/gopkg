//@File     producer.go
//@Time     2023/04/29
//@Author   #Suyghur,

package saramakafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/yyxxgame/gopkg/mq"
	"github.com/yyxxgame/gopkg/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	IProducer interface {
		Publish(topic, key, payload string) error
		PublishCtx(ctx context.Context, topic, key, payload string) error
		Release()
	}

	producer struct {
		*OptionConf
		hooks     []ProducerHook
		finalHook ProducerHook
		sarama.SyncProducer
	}
)

func NewProducer(brokers []string, opts ...Option) IProducer {
	p := &producer{
		OptionConf: &OptionConf{
			producerHooks: []ProducerHook{},
		},
		hooks: []ProducerHook{},
	}
	for _, opt := range opts {
		opt(p.OptionConf)
	}

	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	if p.partitioner == nil {
		p.partitioner = sarama.NewRoundRobinPartitioner
	}
	config.Producer.Partitioner = p.partitioner

	if p.username == "" || p.password == "" {
		config.Net.SASL.Enable = false
	} else {
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		config.Net.SASL.Version = sarama.SASLHandshakeV0
		config.Net.SASL.Handshake = true
		config.Net.SASL.User = p.username
		config.Net.SASL.Password = p.password

		config.Net.TLS.Enable = false
	}

	syncProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logx.Errorf("[SARAMA-KAFKA-ERROR]: MustNewProducer on error: %v", err)
		panic(err)
	}
	p.SyncProducer = syncProducer

	if p.tracer != nil {
		p.hooks = append(p.hooks, newProducerTraceHook(p.tracer).Handle)
	}

	p.hooks = append(p.hooks, newProducerDurationHook().Handle)

	p.hooks = append(p.hooks, p.producerHooks...)

	p.finalHook = chainProducerHooks(p.hooks...)

	return p
}

func (p *producer) Publish(topic, key, payload string) error {
	return p.PublishCtx(context.Background(), topic, key, payload)
}

func (p *producer) PublishCtx(ctx context.Context, topic, key, payload string) error {
	traceId := xtrace.GetTraceId(ctx).String()
	message := &sarama.ProducerMessage{}
	message.Key = sarama.StringEncoder(key)
	message.Topic = topic
	message.Value = sarama.StringEncoder(payload)
	message.Headers = []sarama.RecordHeader{
		{
			Key:   sarama.ByteEncoder(mq.HeaderTraceId),
			Value: sarama.ByteEncoder(traceId),
		},
	}

	return p.finalHook(ctx, message, func(ctx context.Context, message *sarama.ProducerMessage) error {
		return p.produce(ctx, message)
	})
}

func (p *producer) produce(ctx context.Context, message *sarama.ProducerMessage) error {
	partition, offset, err := p.SendMessage(message)
	if err != nil {
		logx.WithContext(ctx).Errorf("[SARAMA-KAFKA-ERROR]: publishMessage.SendMessage to topic: %s, on error: %v", message.Topic, err)
		return err
	}

	logx.WithContext(ctx).Infof("[SARAMA-KAFKA]: publishMessage.SendMessage to topic: %s, on success, partition: %d, offset: %v", message.Topic, partition, offset)
	return nil
}

func (p *producer) Release() {
	_ = p.Close()
}
