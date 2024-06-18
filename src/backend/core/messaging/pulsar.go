package messaging

//
//import (
//	"cognix.ch/api/v2/core/proto"
//	"context"
//	"fmt"
//	"github.com/apache/pulsar-client-go/pulsar"
//	proto2 "github.com/golang/protobuf/proto"
//	"go.opentelemetry.io/otel"
//	"go.opentelemetry.io/otel/propagation"
//	"go.uber.org/zap"
//	"time"
//)
//
//type (
//	pulsarConfig struct {
//		URL               string `env:"PULSAR_URL"`
//		OperationTimeout  int    `env:"OPERATION_TIMEOUT" envDefault:"30"`
//		ConnectionTimeout int    `env:"CONNECTION_TIMEOUT" envDefault:"30"`
//		ReconsumeTimeout  int    `env:"RECONSUME_TIMEOUT" envDefault:"5"`
//	}
//	pulsarClient struct {
//		ReconsumeTimeout time.Duration
//		conn             pulsar.Client
//		producers        map[string]pulsar.Producer
//		subscriber       map[string]pulsar.Consumer
//	}
//)
//
//func (p *pulsarClient) StreamConfig() *StreamConfig {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (p *pulsarClient) Publish(ctx context.Context, topic string, body proto2.Message) error {
//	msg, err := buildMessageAny(ctx, body)
//	if err != nil {
//		return err
//	}
//	producer, ok := p.producers[topic]
//	if !ok {
//		producer, err = p.conn.CreateProducer(pulsar.ProducerOptions{
//			Topic:  topic,
//			Schema: pulsar.NewProtoNativeSchemaWithMessage(&proto.ConnectorRequest{}, nil),
//		})
//		if err != nil {
//			return err
//		}
//		p.producers[topic] = producer
//	}
//
//	_, err = producer.Send(ctx, &pulsar.ProducerMessage{
//		Payload: msg,
//	})
//	return err
//}
//
//func (p *pulsarClient) Listen(ctx context.Context, topic, subscriptionName string, handler MessageHandler) error {
//	consumer, err := p.conn.Subscribe(pulsar.ConsumerOptions{
//		Topic:                          topic,
//		SubscriptionName:               subscriptionName,
//		Type:                           pulsar.Shared,
//		RetryEnable:                    false,
//		NackRedeliveryDelay:            0,
//		MaxPendingChunkedMessage:       0,
//		ExpireTimeOfIncompleteChunk:    0,
//		AutoAckIncompleteChunk:         false,
//		EnableBatchIndexAcknowledgment: false,
//		SubscriptionMode:               0,
//		StartMessageIDInclusive:        false,
//	})
//	if err != nil {
//		return err
//	}
//	defer consumer.Close()
//	for {
//		// may block here
//		select {
//		case <-ctx.Done():
//			break
//		default:
//
//		}
//		msg, err := consumer.Receive(ctx)
//		if err != nil {
//			zap.S().Errorf("Receive message error: %s", err.Error())
//			break
//		}
//		if err = consumer.Ack(msg); err != nil {
//			zap.S().Errorf("Ack message error: %s", err.Error())
//		}
//
//		var message proto.Message
//
//		if err = proto2.Unmarshal(msg.Payload(), &message); err != nil {
//			continue
//		}
//		if err = handler(ctx, &message); err != nil {
//			consumer.ReconsumeLater(msg, p.ReconsumeTimeout)
//			zap.S().Errorf("Reconsume message error: %s", err.Error())
//		}
//	}
//	if err = consumer.Unsubscribe(); err != nil {
//		zap.S().Errorf("Unsubscribe message error: %s", err.Error())
//	}
//	return nil
//}
//
//func (p *pulsarClient) processMessage(ctx context.Context, consumer pulsar.Consumer, handler MessageHandler, msg pulsar.Message) (err error) {
//	if err = consumer.Ack(msg); err != nil {
//		zap.S().Errorf("Ack message error: %s", err.Error())
//	}
//	var message proto.Message
//	if err := proto2.Unmarshal(msg.Payload(), &message); err != nil {
//		return fmt.Errorf("error unmarshalling message: %s", string(msg.Payload()))
//	}
//	return handler(ctx, &message)
//}
//
//func (p *pulsarClient) Close() {
//	for _, producer := range p.producers {
//		producer.Close()
//	}
//	p.conn.Close()
//}
//
//func NewPulsar(cfg *pulsarConfig) (Client, error) {
//	coon, err := pulsar.NewClient(pulsar.ClientOptions{
//		URL:               cfg.URL,
//		ConnectionTimeout: time.Duration(cfg.ConnectionTimeout) * time.Second,
//		OperationTimeout:  time.Duration(cfg.OperationTimeout) * time.Second,
//	})
//	if err != nil {
//		return nil, err
//	}
//	return &pulsarClient{
//		conn:             coon,
//		ReconsumeTimeout: time.Second * time.Duration(cfg.ReconsumeTimeout),
//		producers:        make(map[string]pulsar.Producer),
//		subscriber:       make(map[string]pulsar.Consumer),
//	}, nil
//}
//
//func buildMessageAny(ctx context.Context, body proto2.Message) ([]byte, error) {
//	return proto2.Marshal(body)
//}
//
//func buildMessage(ctx context.Context, body *proto.Body) ([]byte, error) {
//	header := make(propagation.MapCarrier)
//	otel.GetTextMapPropagator().Inject(ctx, &header)
//	msg := &proto.Message{
//		Header: header,
//		Body:   body,
//	}
//	return proto2.Marshal(msg)
//}
