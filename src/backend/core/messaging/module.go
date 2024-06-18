package messaging

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	proto2 "github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

const (
	providerNats   = "nats"
	providerPulsar = "pulsar"
)

type (
	Config struct {
		Nats *natsConfig
		//Pulsar *pulsarConfig
		Stream *StreamConfig
	}
	natsConfig struct {
		URL                  string `env:"NATS_CLIENT_URL,required"`
		ConnectTimeout       int    `env:"NATS_CLIENT_CONNECT_TIMEOUT" envDefault:"3"`
		ReconnectTimeout     int    `env:"NATS_CLIENT_RECONNECT_TIME_WAIT" envDefault:"30"`
		MaxReconnectAttempts int    `env:"NATS_CLIENT_MAX_RECONNECT_ATTEMPTS" envDefault:"3"`
	}
	// StreamConfig contains variables for configure streams
	StreamConfig struct {
		ConnectorStreamName    string `env:"NATS_CLIENT_CONNECTOR_STREAM_NAME,required"`
		ConnectorStreamSubject string `env:"NATS_CLIENT_CONNECTOR_STREAM_SUBJECT,required"`
		SemanticStreamName     string `env:"NATS_CLIENT_SEMANTIC_STREAM_NAME,required"`
		SemanticStreamSubject  string `env:"NATS_CLIENT_SEMANTIC_STREAM_SUBJECT,required"`
		AckWait                int    `env:"NATS_CLIENT_CONNECTOR_ACK_WAIT,required"`
		MaxDeliver             int    `env:"NATS_CLIENT_CONNECTOR_MAX_DELIVER,required"`
	}

	MessageHandler func(ctx context.Context, msg jetstream.Msg) error
	Client         interface {
		Publish(ctx context.Context, streamName, topic string, body proto2.Message) error
		Listen(ctx context.Context, streamName, topic string, handler MessageHandler) error
		StreamConfig() *StreamConfig
		Close()
	}
)

const (
	reconnectAttempts = 120
	reconnectWaitTime = 5 * time.Second
	streamMaxPending  = 256
)

var NatsModule = fx.Options(
	fx.Provide(func() (*Config, error) {
		cfg := Config{
			Nats:   &natsConfig{},
			Stream: &StreamConfig{},
		}
		err := utils.ReadConfig(&cfg)
		if err != nil {
			zap.S().Errorf(err.Error())
			return nil, err
		}
		return &cfg, nil
	},
		NewClientStream,
	),
)
