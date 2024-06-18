package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type (
	TenantRepository interface {
		GetUsers(ctx context.Context, tenantID uuid.UUID) ([]*model.User, error)
	}
	tenantRepository struct {
		db *pg.DB
	}
)

func (r *tenantRepository) GetUsers(ctx context.Context, tenantID uuid.UUID) ([]*model.User, error) {
	users := make([]*model.User, 0)
	if err := r.db.Model(&users).Where("tenant_id = ?", tenantID).Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "cannot get users")
	}
	return users, nil
}

func NewTenantRepository(db *pg.DB) TenantRepository {
	return &tenantRepository{db: db}
}
