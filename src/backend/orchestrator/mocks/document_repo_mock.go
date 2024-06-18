package mocks

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"context"
	"github.com/shopspring/decimal"
)

type MockDocumentRepository struct{}

func (m MockDocumentRepository) FindByConnectorIDAndUser(ctx context.Context, user *model.User, connectorID int64) ([]*model.Document, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockDocumentRepository) FindByConnectorID(ctx context.Context, connectorID int64) ([]*model.Document, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockDocumentRepository) FindByID(ctx context.Context, id int64) (*model.Document, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockDocumentRepository) Create(ctx context.Context, document *model.Document) error {
	document.ID = decimal.NewFromInt(1)
	return nil
}

func (m MockDocumentRepository) Update(ctx context.Context, document *model.Document) error {
	return nil
}

func (m MockDocumentRepository) DeleteByIDS(ctx context.Context, ids ...int64) error {
	//TODO implement me
	panic("implement me")
}

func NewMockDocumentRepo() repository.DocumentRepository {
	return &MockDocumentRepository{}
}
