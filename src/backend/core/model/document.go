package model

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Document struct {
	tableName       struct{}            `pg:"documents"`
	ID              decimal.Decimal     `json:"id,omitempty"`
	ParentID        decimal.NullDecimal `json:"parent_id,omitempty" pg:",use_zero"`
	SourceID        string              `json:"source_id,omitempty"`
	ConnectorID     decimal.Decimal     `json:"connector_id,omitempty"`
	URL             string              `json:"url,omitempty" pg:"url"`
	Signature       string              `json:"signature,omitempty" pg:",use_zero"`
	ChunkingSession uuid.NullUUID       `json:"chunking_session,omitempty" pg:",use_zero"`
	Analyzed        bool                `json:"analyzed" pg:",use_zero"`
	CreationDate    time.Time           `json:"creation_date,omitempty"`
	LastUpdate      pg.NullTime         `json:"last_update,omitempty" pg:",use_zero"`
	IsExists        bool                `json:"-" pg:"-"`
}

type DocumentResponse struct {
	ID          decimal.Decimal `json:"id,omitempty"`
	MessageID   decimal.Decimal `json:"message_id,omitempty"`
	Link        string          `json:"link,omitempty"`
	DocumentID  string          `json:"document_id,omitempty"`
	Content     string          `json:"content,omitempty"`
	UpdatedDate time.Time       `json:"updated_date,omitempty"`
}
