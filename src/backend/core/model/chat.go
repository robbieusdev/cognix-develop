package model

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

const (
	MessageTypeUser      = "user"
	MessageTypeAssistant = "assistant"
	MessageTypeSystem    = "system"
)

type (
	ChatSession struct {
		tableName    struct{}        `pg:"chat_sessions"`
		ID           decimal.Decimal `json:"id,omitempty"`
		UserID       uuid.UUID       `json:"user_id,omitempty"`
		Description  string          `json:"description,omitempty"`
		CreationDate time.Time       `json:"creation_date,omitempty"`
		DeletedDate  pg.NullTime     `json:"deleted_date,omitempty"`
		PersonaID    decimal.Decimal `json:"persona_id,omitempty"`
		OneShot      bool            `json:"one_shot,omitempty" pg:",use_zero"`
		Messages     []*ChatMessage  `json:"messages,omitempty" pg:"rel:has-many"`
		Persona      *Persona        `json:"persona,omitempty" pg:"rel:has-one"`
	}

	ChatMessage struct {
		tableName          struct{}                   `pg:"chat_messages"`
		ID                 decimal.Decimal            `json:"id,omitempty"`
		ChatSessionID      decimal.Decimal            `json:"chat_session_id,omitempty"`
		Message            string                     `json:"message,omitempty"  pg:",use_zero"`
		MessageType        string                     `json:"message_type,omitempty"`
		TimeSent           time.Time                  `json:"time_sent,omitempty"`
		TokenCount         int                        `json:"token_count,omitempty" pg:",use_zero"`
		ParentMessageID    decimal.Decimal            `json:"parent_message,omitempty" pg:"parent_message,use_zero"`
		LatestChildMessage int                        `json:"latest_child_message,omitempty" pg:",use_zero"`
		RephrasedQuery     string                     `json:"rephrased_query,omitempty" pg:",use_zero"`
		Citations          []*DocumentResponse        `json:"citations,omitempty" pg:"-"`
		Error              string                     `json:"error,omitempty" pg:",use_zero"`
		Feedback           *ChatMessageFeedback       `json:"feedback,omitempty" pg:"rel:has-one,fk:id,join_fk:chat_message_id"`
		ParentMessage      *ChatMessage               `json:"-" pg:"-"`
		DocumentPairs      []*ChatMessageDocumentPair `json:"-" pg:"rel:has-many,fk:chat_message_id,join_fk:chat_message_id"`
	}

	ChatMessageFeedback struct {
		tableName     struct{}        `pg:"chat_message_feedbacks"`
		ID            decimal.Decimal `json:"id,omitempty"`
		ChatMessageID decimal.Decimal `json:"chat_message_id,omitempty"`
		UserID        uuid.UUID       `json:"user_id,omitempty"`
		UpVotes       bool            `json:"up_votes" pg:",use_zero"`
		Feedback      string          `json:"feedback,omitempty" pg:",use_zero"`
	}
	// ChatMessageDocumentPair struct represent data from table chat_message_document_pairs
	// what documents were found for each message.
	ChatMessageDocumentPair struct {
		tableName     struct{}        `pg:"chat_message_document_pairs"`
		ID            decimal.Decimal `json:"id"`
		ChatMessageID decimal.Decimal `json:"chat_message_id" pg:",use_zero"`
		DocumentID    decimal.Decimal `json:"document_id" pg:",use_zero"`
		Document      *Document       `json:"document" pg:"rel:has-one"`
	}
)

var _ pg.AfterSelectHook = (*ChatMessage)(nil)

func (c *ChatMessage) AfterSelect(ctx context.Context) error {
	for _, dp := range c.DocumentPairs {
		if dp.Document == nil {
			continue
		}
		doc := &DocumentResponse{
			ID:         dp.DocumentID,
			MessageID:  dp.ChatMessageID,
			Link:       dp.Document.URL,
			DocumentID: dp.Document.SourceID,
		}
		if !dp.Document.LastUpdate.IsZero() {
			doc.UpdatedDate = dp.Document.LastUpdate.Time
		} else {
			doc.UpdatedDate = dp.Document.CreationDate
		}
		c.Citations = append(c.Citations, doc)
	}
	return nil
}
