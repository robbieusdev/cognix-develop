package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/google/uuid"
)

type ChatRepository interface {
	GetSessions(ctx context.Context, userID uuid.UUID) ([]*model.ChatSession, error)
	GetSessionByID(ctx context.Context, userID uuid.UUID, id int64) (*model.ChatSession, error)
	CreateSession(ctx context.Context, session *model.ChatSession) error
	SendMessage(ctx context.Context, message *model.ChatMessage) error
	UpdateMessage(ctx context.Context, message *model.ChatMessage) error
	GetMessageByIDAndUserID(ctx context.Context, id int64, userID uuid.UUID) (*model.ChatMessage, error)
	MessageFeedback(ctx context.Context, feedback *model.ChatMessageFeedback) error
}

type chatRepository struct {
	db *pg.DB
}

func (r *chatRepository) GetMessageByIDAndUserID(ctx context.Context, id int64, userID uuid.UUID) (*model.ChatMessage, error) {
	var message model.ChatMessage
	if err := r.db.Model(&message).
		Relation("Feedback").
		Join("inner join chat_sessions on chat_sessions.id = chat_message.chat_session_id and chat_sessions.user_id = ?", userID).
		Where("chat_message.id = ?", id).First(); err != nil {
		return nil, utils.NotFound.Wrap(err, "cannot find message by id")
	}
	return &message, nil
}

func (r *chatRepository) MessageFeedback(ctx context.Context, feedback *model.ChatMessageFeedback) error {
	stm := r.db.WithContext(ctx).Model(feedback)
	if feedback.ID.IntPart() == 0 {
		if _, err := stm.Insert(); err != nil {
			return utils.Internal.Wrap(err, "can not add feedback")
		}
		return nil
	}
	if _, err := stm.Where("id = ?", feedback.ID).Update(); err != nil {
		return utils.Internal.Wrap(err, "can not update feedback")
	}
	return nil
}

func (r *chatRepository) SendMessage(ctx context.Context, message *model.ChatMessage) error {
	if _, err := r.db.WithContext(ctx).Model(message).Insert(); err != nil {
		return utils.Internal.Wrap(err, "can not save message")
	}
	return nil
}
func (r *chatRepository) UpdateMessage(ctx context.Context, message *model.ChatMessage) error {
	return r.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		if _, err := tx.Model(message).Where("id = ?", message.ID).Update(); err != nil {
			return utils.Internal.Wrap(err, "can not save message")
		}
		if len(message.DocumentPairs) > 0 {
			if _, err := tx.Model(&message.DocumentPairs).Insert(); err != nil {
				return utils.Internal.Wrap(err, "can not create document pairs")
			}
		}
		return nil
	})

}

func NewChatRepository(db *pg.DB) ChatRepository {
	return &chatRepository{db: db}
}
func (r *chatRepository) GetSessions(ctx context.Context, userID uuid.UUID) ([]*model.ChatSession, error) {
	sessions := make([]*model.ChatSession, 0)
	if err := r.db.WithContext(ctx).Model(&sessions).
		Where("user_id = ?", userID).
		Where("deleted_date is null").
		Order("creation_date desc").Select(); err != nil {
		return nil, utils.NotFound.Wrapf(err, "can not find sessions")
	}
	return sessions, nil
}

func (r *chatRepository) GetSessionByID(ctx context.Context, userID uuid.UUID, id int64) (*model.ChatSession, error) {
	var session model.ChatSession
	if err := r.db.WithContext(ctx).Model(&session).
		Where("chat_session.user_id = ?", userID).
		Where("chat_session.id = ?", id).
		Relation("Persona").
		Relation("Persona.Prompt").
		Relation("Persona.LLM").
		Relation("Messages", func(query *orm.Query) (*orm.Query, error) {
			return query.Order("time_sent asc"), nil
		}).
		Relation("Messages.Feedback").
		Relation("Messages.DocumentPairs").
		Relation("Messages.DocumentPairs.Document").
		First(); err != nil {
		return nil, utils.NotFound.Wrapf(err, "can not find session")
	}
	for _, msg := range session.Messages {
		_ = msg.AfterSelect(ctx)
	}
	return &session, nil
}

func (r *chatRepository) CreateSession(ctx context.Context, session *model.ChatSession) error {
	if _, err := r.db.WithContext(ctx).Model(session).Insert(); err != nil {
		return utils.Internal.Wrap(err, "can not create chat session")
	}
	return nil
}
