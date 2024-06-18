package mocks

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"context"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"time"
)

var MockedConnectors = map[int64]*model.Connector{
	1: {
		ID:   decimal.NewFromInt(1),
		Name: "web connector ready to process",
		Type: model.SourceTypeWEB,
		ConnectorSpecificConfig: model.JSONMap{
			"url": "http://testurl",
		},
		RefreshFreq:       60,
		Status:            model.ConnectorStatusReadyToProcessed,
		TotalDocsAnalyzed: 0,
		CreationDate:      time.Now().UTC(),
		User: &model.User{
			EmbeddingModel: &model.EmbeddingModel{
				ModelID:  "",
				ModelDim: 3,
			},
		},
	},
	2: {
		ID:   decimal.NewFromInt(2),
		Name: "File connector ready to process",
		Type: model.SourceTypeFile,
		ConnectorSpecificConfig: model.JSONMap{
			"file_name": "file name",
		},
		RefreshFreq:       60,
		Status:            model.ConnectorStatusReadyToProcessed,
		TotalDocsAnalyzed: 0,
		CreationDate:      time.Now().UTC(),
		User: &model.User{
			EmbeddingModel: &model.EmbeddingModel{
				ModelID:  "",
				ModelDim: 3,
			},
		},
	},
	3: {
		ID:   decimal.NewFromInt(3),
		Name: "One drive connector disabled",
		Type: model.SourceTypeOneDrive,
		ConnectorSpecificConfig: model.JSONMap{
			"file_name": "file name",
		},
		RefreshFreq:       60,
		Status:            model.ConnectorStatusDisabled,
		TotalDocsAnalyzed: 0,
		CreationDate:      time.Now().UTC(),
		User: &model.User{
			EmbeddingModel: &model.EmbeddingModel{
				ModelID:  "",
				ModelDim: 3,
			},
		},
	},
	4: {
		ID:   decimal.NewFromInt(2),
		Name: "File connector without embedding model",
		Type: model.SourceTypeFile,
		ConnectorSpecificConfig: model.JSONMap{
			"file_name": "file name",
		},
		RefreshFreq:       60,
		Status:            model.ConnectorStatusReadyToProcessed,
		TotalDocsAnalyzed: 0,
		CreationDate:      time.Now().UTC(),
		User:              &model.User{},
	},
	5: {
		ID:   decimal.NewFromInt(5),
		Name: "One drive connector ready for processing ",
		Type: model.SourceTypeOneDrive,
		ConnectorSpecificConfig: model.JSONMap{
			"file_name": "file name",
		},
		RefreshFreq:       60,
		Status:            model.ConnectorStatusReadyToProcessed,
		TotalDocsAnalyzed: 0,
		CreationDate:      time.Now().UTC(),
		User: &model.User{
			EmbeddingModel: &model.EmbeddingModel{
				ModelID:  "",
				ModelDim: 3,
			},
		},
	},
}

type MockConnectorRepo struct {
	workCh            chan int
	expectedIteration int
	iteration         int
}

func (m *MockConnectorRepo) GetActive(ctx context.Context) ([]*model.Connector, error) {
	result := make([]*model.Connector, 0)
	zap.S().Errorf("load iteration %d", m.iteration)
	if m.iteration >= m.expectedIteration {
		close(m.workCh)
		return result, nil
	}
	m.iteration++
	zap.S().Info("connector for running ")
	for _, conn := range MockedConnectors {
		if conn.Status == model.ConnectorStatusReadyToProcessed ||
			conn.Status == model.ConnectorStatusSuccess ||
			conn.Status == model.ConnectorStatusError {
			zap.S().Infof("| %s \t\t| %s \t\t | %s \t\t | %v |",
				conn.Name, conn.Type, conn.Status, conn.LastUpdate)
			result = append(result, conn)
		}
	}
	return result, nil
}

func (m *MockConnectorRepo) GetAllByUser(ctx context.Context, tenantID, userID uuid.UUID) ([]*model.Connector, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockConnectorRepo) GetByIDAndUser(ctx context.Context, tenantID, userID uuid.UUID, id int64) (*model.Connector, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockConnectorRepo) GetByID(ctx context.Context, id int64) (*model.Connector, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockConnectorRepo) GetBySource(ctx context.Context, tenantID, userID uuid.UUID, source model.SourceType) (*model.Connector, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockConnectorRepo) Create(ctx context.Context, connector *model.Connector) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockConnectorRepo) Update(ctx context.Context, connector *model.Connector) error {
	MockedConnectors[connector.ID.IntPart()] = connector
	return nil
}

func NewMockConnectorRepo(iteration int, workCh chan int) repository.ConnectorRepository {
	return &MockConnectorRepo{
		workCh:            workCh,
		expectedIteration: iteration,
	}
}
