package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Prompt struct {
	tableName    struct{}        `pg:"prompts"`
	ID           decimal.Decimal `json:"id,omitempty"`
	PersonaID    decimal.Decimal `json:"persona_id,omitempty"`
	UserID       uuid.UUID       `json:"user_id,omitempty"`
	Name         string          `json:"name,omitempty"`
	Description  string          `json:"description,omitempty" pg:",use_zero"`
	SystemPrompt string          `json:"system_prompt,omitempty" pg:",use_zero"`
	TaskPrompt   string          `json:"task_prompt,omitempty" pg:",use_zero"`
	CreationDate time.Time       `json:"creation_date,omitempty"`
	LastUpdate   pg.NullTime     `json:"last_update,omitempty" pg:",use_zero"`
	DeletedDate  pg.NullTime     `json:"deleted_date,omitempty" pg:",use_zero"`
}
