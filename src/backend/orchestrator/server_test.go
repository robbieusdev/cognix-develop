package main

import (
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/utils"
	"cognix.ch/api/v2/orchestrator/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestOrchestrator_Scheduler(t *testing.T) {
	//b := &bytes.Buffer{}
	// call the constructor from your test code with the arbitrary writer
	//mycore := NewCustomLogger(os.Stderr)
	//zap.ReplaceGlobals(zap.New(mycore))
	utils.InitLogger(true)
	workCh := make(chan int, 10)

	srv, err := NewServer(
		&Config{
			RenewInterval: 15,
			FileSizeLimit: 1,
		},
		mocks.NewMockConnectorRepo(10, workCh),
		mocks.NewMockDocumentRepo(),
		mocks.NewMockMessenger(workCh),
		&messaging.Config{
			Stream: &messaging.StreamConfig{
				ConnectorStreamName: "connector",
				SemanticStreamName:  "semantic",
			},
		},
	)
	zap.S().Infof("connecto in database ")
	for _, conn := range mocks.MockedConnectors {
		zap.S().Infof("| %s \t\t| %s \t\t | %s \t\t | %v |",
			conn.Name, conn.Type, conn.Status, conn.LastUpdate)
	}

	assert.NoError(t, err)
	srv.run(context.Background())
	for i := range workCh {
		t.Logf("iteration %d", i)
	}
	//t.Log(b.String())
}
