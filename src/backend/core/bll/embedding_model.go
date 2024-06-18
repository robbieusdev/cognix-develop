package bll

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/repository"
	"context"
	"github.com/go-pg/pg/v10"
	"time"
)

type (
	EmbeddingModelBL interface {
		GetAll(ctx context.Context, user *model.User, param *parameters.ArchivedParam) ([]*model.EmbeddingModel, error)
		GetByID(ctx context.Context, user *model.User, id int64) (*model.EmbeddingModel, error)
		Create(ctx context.Context, user *model.User, em *parameters.EmbeddingModelParam) (*model.EmbeddingModel, error)
		Update(ctx context.Context, user *model.User, id int64, em *parameters.EmbeddingModelParam) (*model.EmbeddingModel, error)
		Delete(ctx context.Context, user *model.User, id int64) (*model.EmbeddingModel, error)
		Restore(ctx context.Context, user *model.User, id int64) (*model.EmbeddingModel, error)
	}
	embeddingModelBL struct {
		emRepo repository.EmbeddingModelRepository
	}
)

func (b *embeddingModelBL) Create(ctx context.Context, user *model.User, em *parameters.EmbeddingModelParam) (*model.EmbeddingModel, error) {
	embeddingModel := model.EmbeddingModel{
		TenantID:     user.TenantID,
		ModelID:      em.ModelID,
		ModelName:    em.ModelName,
		ModelDim:     em.ModelDim,
		URL:          em.URL,
		IsActive:     em.IsActive,
		CreationDate: time.Now().UTC(),
	}
	if err := b.emRepo.Create(ctx, &embeddingModel); err != nil {
		return nil, err
	}
	return &embeddingModel, nil
}

func (b *embeddingModelBL) Update(ctx context.Context, user *model.User, id int64, em *parameters.EmbeddingModelParam) (*model.EmbeddingModel, error) {
	existingEM, err := b.emRepo.GetByID(ctx, user.TenantID, id)
	if err != nil {
		return nil, err
	}
	existingEM.ModelID = em.ModelID
	existingEM.ModelName = em.ModelName
	existingEM.ModelDim = em.ModelDim
	existingEM.URL = em.URL
	existingEM.IsActive = em.IsActive
	existingEM.LastUpdate = pg.NullTime{time.Now().UTC()}
	if err = b.emRepo.Update(ctx, existingEM); err != nil {
		return nil, err
	}
	return existingEM, nil
}

func (b *embeddingModelBL) Delete(ctx context.Context, user *model.User, id int64) (*model.EmbeddingModel, error) {
	existingEM, err := b.emRepo.GetByID(ctx, user.TenantID, id)
	if err != nil {
		return nil, err
	}
	existingEM.DeletedDate = pg.NullTime{time.Now().UTC()}
	if err = b.emRepo.Update(ctx, existingEM); err != nil {
		return nil, err
	}
	return existingEM, nil
}

func (b *embeddingModelBL) Restore(ctx context.Context, user *model.User, id int64) (*model.EmbeddingModel, error) {
	existingEM, err := b.emRepo.GetByID(ctx, user.TenantID, id)
	if err != nil {
		return nil, err
	}
	existingEM.DeletedDate = pg.NullTime{time.Time{}}
	if err = b.emRepo.Update(ctx, existingEM); err != nil {
		return nil, err
	}
	return existingEM, nil
}

func (b *embeddingModelBL) GetAll(ctx context.Context, user *model.User, param *parameters.ArchivedParam) ([]*model.EmbeddingModel, error) {
	return b.emRepo.GetAll(ctx, user.TenantID, param)
}

func (b *embeddingModelBL) GetByID(ctx context.Context, user *model.User, id int64) (*model.EmbeddingModel, error) {
	return b.emRepo.GetByID(ctx, user.TenantID, id)
}

func NewEmbeddingModelBL(emRepo repository.EmbeddingModelRepository) EmbeddingModelBL {
	return &embeddingModelBL{emRepo: emRepo}
}
