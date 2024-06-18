package messaging

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClientStream_Listen(t *testing.T) {
	cfg := &natsConfig{
		URL:                 "localhost:4222",
		ConnectorStreamName: "test-1",
	}
	stClient, err := NewClientStream(cfg)
	assert.Nil(t, err)

	go func() {
		time.Sleep(time.Second * time.Duration(30))
		t.Error("close client stream")
		stClient.Close()
	}()
	stClient.Listen(context.Background(), model.TopicExecutor, model.SubscriptionExecutor, func(ctx context.Context, msg *proto.Message) error {
		t.Log(msg.Body.String())
		return nil
	})
}
