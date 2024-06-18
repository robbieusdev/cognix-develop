package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type (
	LLMRepository interface {
		GetAll(ctx context.Context) ([]*model.LLM, error)
		GetByUserID(ctx context.Context, userID uuid.UUID) (*model.LLM, error)
	}
	llmRepository struct {
		db *pg.DB
	}
)

func (r *llmRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.LLM, error) {
	var llm model.LLM
	if err := r.db.WithContext(ctx).Model(&llm).
		Join("inner join users on llm.tenant_id = users.tenant_id").
		Where("users.id = ?", userID).
		Limit(1).First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find llm ")
	}
	return &llm, nil
}

func NewLLMRepository(db *pg.DB) LLMRepository {
	return &llmRepository{db: db}
}

func (r *llmRepository) GetAll(ctx context.Context) ([]*model.LLM, error) {
	llms := make([]*model.LLM, 0)
	if err := r.db.WithContext(ctx).Model(&llms).Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find llm")
	}
	return llms, nil
}
