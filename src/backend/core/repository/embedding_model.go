package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"time"
)

type (
	EmbeddingModelRepository interface {
		GetAll(ctx context.Context, tenantID uuid.UUID, param *parameters.ArchivedParam) ([]*model.EmbeddingModel, error)
		GetByID(ctx context.Context, tenantID uuid.UUID, id int64) (*model.EmbeddingModel, error)
		GetDefault(ctx context.Context, tenantID uuid.UUID) (*model.EmbeddingModel, error)
		Create(ctx context.Context, em *model.EmbeddingModel) error
		Update(ctx context.Context, em *model.EmbeddingModel) error
		Delete(ctx context.Context, em *model.EmbeddingModel) error
	}
	embeddingModelRepository struct {
		db *pg.DB
	}
)

func (r *embeddingModelRepository) GetDefault(ctx context.Context, tenantID uuid.UUID) (*model.EmbeddingModel, error) {
	var em model.EmbeddingModel
	if err := r.db.Model(&em).Where("tenant_id = ?", tenantID).
		Where("is_active = true").
		Limit(1).
		First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "Cannot get default embedding model")
	}
	return &em, nil
}

func (r *embeddingModelRepository) GetAll(ctx context.Context, tenantID uuid.UUID, param *parameters.ArchivedParam) ([]*model.EmbeddingModel, error) {
	ems := make([]*model.EmbeddingModel, 0)
	stm := r.db.WithContext(ctx).Model(&ems).Where("tenant_id = ?", tenantID)
	if !param.Archived {
		stm = stm.Where("deleted_date is null")
	}

	if err := stm.Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find embedding models")
	}
	return ems, nil
}

func (r *embeddingModelRepository) GetByID(ctx context.Context, tenantID uuid.UUID, id int64) (*model.EmbeddingModel, error) {
	var em model.EmbeddingModel
	if err := r.db.WithContext(ctx).Model(&em).Where("id = ?", id).
		Where("tenant_id = ?", tenantID).
		Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find embedding models")
	}
	return &em, nil
}

func (r *embeddingModelRepository) Create(ctx context.Context, em *model.EmbeddingModel) error {
	if _, err := r.db.WithContext(ctx).Model(em).Insert(); err != nil {
		return utils.Internal.Wrap(err, "can not create embedding models")
	}
	return nil
}

func (r *embeddingModelRepository) Update(ctx context.Context, em *model.EmbeddingModel) error {
	em.LastUpdate = pg.NullTime{time.Now().UTC()}
	if _, err := r.db.WithContext(ctx).Model(em).
		Where("id = ?", em.ID).
		Update(); err != nil {
		return utils.Internal.Wrap(err, "can not update embedding models")
	}
	return nil
}

func (r *embeddingModelRepository) Delete(ctx context.Context, em *model.EmbeddingModel) error {
	em.DeletedDate = pg.NullTime{time.Now().UTC()}
	if _, err := r.db.WithContext(ctx).Model(em).
		Where("id = ?", em.ID).
		Update(); err != nil {
		return utils.Internal.Wrap(err, "can not update embedding models")
	}
	return nil
}

func NewEmbeddingModelRepository(db *pg.DB) EmbeddingModelRepository {
	return &embeddingModelRepository{
		db: db,
	}
}
