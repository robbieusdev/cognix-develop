package bll

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/storage"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"time"
)

type (
	DocumentBL interface {
		UploadDocument(ctx context.Context, user *model.User, fileName, contentType string, file io.Reader) (*model.Document, error)
	}
	documentBL struct {
		documentRepo  repository.DocumentRepository
		minioClient   storage.MinIOClient
		connectorRepo repository.ConnectorRepository
	}
)

func (b *documentBL) UploadDocument(ctx context.Context, user *model.User, fileName, contentType string, file io.Reader) (*model.Document, error) {

	fileURL, signature, err := b.minioClient.Upload(ctx, model.BucketName(user.TenantID),
		fmt.Sprintf("user-%s/%s-%s", user.ID.String(), uuid.New().String(), fileName), contentType, file)
	if err != nil {
		return nil, err
	}
	connector, err := b.connectorRepo.GetBySource(ctx, user.TenantID, user.ID, model.SourceTypeFile)
	if err != nil {
		return nil, err
	}
	document := &model.Document{
		SourceID:     fileURL,
		ConnectorID:  connector.ID,
		URL:          fileURL,
		Signature:    signature,
		CreationDate: time.Now().UTC(),
	}
	if err = b.documentRepo.Create(ctx, document); err != nil {
		return nil, err
	}
	return document, nil
}

func NewDocumentBL(documentRepo repository.DocumentRepository,
	connectorRepo repository.ConnectorRepository,
	minioClient storage.MinIOClient) DocumentBL {
	return &documentBL{documentRepo: documentRepo,
		connectorRepo: connectorRepo,
		minioClient:   minioClient,
	}
}
