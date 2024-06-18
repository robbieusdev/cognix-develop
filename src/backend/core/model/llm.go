package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type LLM struct {
	tableName    struct{}        `pg:"llms"`
	ID           decimal.Decimal `json:"id,omitempty"`
	Name         string          `json:"name,omitempty"`
	ModelID      string          `json:"model_id,omitempty"`
	TenantID     uuid.UUID       `json:"tenant_id,omitempty"`
	Url          string          `json:"url,omitempty"  pg:",use_zero"`
	ApiKey       string          `json:"api_key"`
	Endpoint     string          `json:"endpoint,omitempty"`
	CreationDate time.Time       `json:"creation_date,omitempty"`
	LastUpdate   pg.NullTime     `json:"last_update,omitempty" pg:",use_zero"`

	DeletedDate pg.NullTime `json:"deleted_date,omitempty" pg:",use_zero"`
}

func (l *LLM) MaskApiKey() string {
	if len(l.ApiKey) < 10 {
		return "***"
	}
	return string(l.ApiKey[:4]) + "***" + l.ApiKey[len(l.ApiKey)-4:]
}
