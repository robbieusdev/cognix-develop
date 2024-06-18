package mocks

import (
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"context"
	"github.com/go-pg/pg/v10"
	proto2 "github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"time"
)

type MockMessenger struct {
	workCh chan int
}

func (m MockMessenger) Publish(ctx context.Context, streamName, topic string, body proto2.Message) error {

	if streamName == "semantic" {
		semantic := body.(*proto.SemanticData)
		conn := MockedConnectors[semantic.ConnectorId]
		zap.S().Infof("befor sending to semantic .... ")
		zap.S().Infof("| %s \t\t| %s \t\t | %s \t\t | %v |",
			conn.Name, conn.Type, conn.Status, conn.LastUpdate)

		conn.Status = model.ConnectorStatusSuccess
		conn.LastUpdate = pg.NullTime{time.Now().UTC()}
		zap.S().Infof("emulate semantic  work .... ")
		zap.S().Infof("| %s \t\t| %s \t\t | %s \t\t | %v |",
			conn.Name, conn.Type, conn.Status, conn.LastUpdate)
		return nil
	}
	if streamName == "connector" {
		connRequest := body.(*proto.ConnectorRequest)
		conn := MockedConnectors[connRequest.Id]
		zap.S().Infof("befor sending to connector  .... ")
		zap.S().Infof("| %s \t\t| %s \t\t | %s \t\t | %v |",
			conn.Name, conn.Type, conn.Status, conn.LastUpdate)

		conn.Status = model.ConnectorStatusWorking
		conn.LastUpdate = pg.NullTime{time.Now().UTC()}
		zap.S().Infof("emulate connector  work .... ")
		zap.S().Infof("| %s \t\t| %s \t\t | %s \t\t | %v |",
			conn.Name, conn.Type, conn.Status, conn.LastUpdate)
		return nil
	}
	return nil
}

func (m MockMessenger) Listen(ctx context.Context, streamName, topic string, handler messaging.MessageHandler) error {

	return nil
}

func (m MockMessenger) StreamConfig() *messaging.StreamConfig {
	return &messaging.StreamConfig{
		ConnectorStreamName: "connector",
		SemanticStreamName:  "semantic",
	}
}

func (m MockMessenger) Close() {
	//TODO implement me
	panic("implement me")
}

func NewMockMessenger(workCh chan int) messaging.Client {
	return &MockMessenger{
		workCh: workCh,
	}
}
