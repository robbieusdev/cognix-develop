package messaging

import (
	"context"
	proto2 "github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	_ "github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"sync"
	"time"
)

type clientStream struct {
	js        jetstream.JetStream
	cancel    context.CancelFunc
	ctx       context.Context
	wg        *sync.WaitGroup
	streamCfg *StreamConfig
	ackWait   time.Duration
}

func (c *clientStream) StreamConfig() *StreamConfig {

	return c.streamCfg
}

func (c *clientStream) Close() {
	c.wg.Add(1)
	c.cancel()
	c.wg.Wait()
}

func (c *clientStream) Publish(ctx context.Context, streamName, topic string, body proto2.Message) error {
	_, err := c.js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:      streamName,
		Retention: jetstream.WorkQueuePolicy,
		//AllowDirect: true,
		Subjects: []string{topic},
	})
	if err != nil {
		zap.S().Errorf("Error creating stream: %s", err.Error())
		return err
	}

	message, err := proto2.Marshal(body)
	if err != nil {
		return err
	}
	_, err = c.js.Publish(ctx, topic, message)
	//,
	//		nats.AckWait(time.Minute*2)
	if err != nil {
		return err
	}
	return nil
}

func (c *clientStream) Listen(ctx context.Context, streamName, topic string, handler MessageHandler) error {

	stream, err := c.js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:      streamName,
		Retention: jetstream.WorkQueuePolicy,
		//AllowDirect: true,
		Subjects: []string{topic},
	})
	if err != nil {
		zap.S().Errorf("Error creating stream: %s", err.Error())
		return err
	}

	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       streamName,
		MaxDeliver:    c.streamCfg.MaxDeliver,
		FilterSubject: topic,
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       c.ackWait,
		DeliverPolicy: jetstream.DeliverAllPolicy,
	})
	if err != nil {
		zap.S().Errorf("Failed to create consumer for subscription %v", err)
	}
	cons.Consume(func(msg jetstream.Msg) {
		msg.InProgress()
		if err := handler(ctx, msg); err != nil {
			zap.S().Errorf("Error handling message: %s", err.Error())
		}
		err := msg.Ack()
		if err != nil {
			zap.S().Errorf("Error acknowledging message: %s", err.Error())
		}

	})
	<-c.ctx.Done()
	c.wg.Done()
	return nil
}

func NewClientStream(cfg *Config) (Client, error) {
	conn, err := nats.Connect(
		cfg.Nats.URL,
		nats.Timeout(time.Duration(cfg.Nats.ConnectTimeout)*time.Second),
		nats.ReconnectWait(time.Duration(cfg.Nats.ReconnectTimeout)*time.Second),
		nats.MaxReconnects(cfg.Nats.MaxReconnectAttempts),
	)
	if err != nil {
		zap.S().Errorf("Error connecting to NATS: %s", err.Error())
		return nil, err
	}

	js, err := jetstream.New(conn)
	if err != nil {
		zap.S().Errorf("Error connecting to NATS: %s", err.Error())
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &clientStream{
		js:        js,
		ctx:       ctx,
		cancel:    cancel,
		wg:        &sync.WaitGroup{},
		streamCfg: cfg.Stream,
		ackWait:   time.Duration(cfg.Stream.AckWait) * time.Second,
	}, nil
}
