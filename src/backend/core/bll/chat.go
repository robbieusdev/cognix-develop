package bll

import (
	"cognix.ch/api/v2/core/ai"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/responder"
	"cognix.ch/api/v2/core/storage"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type ChatBL interface {
	GetSessions(ctx context.Context, user *model.User) ([]*model.ChatSession, error)
	GetSessionByID(ctx context.Context, user *model.User, id int64) (*model.ChatSession, error)
	CreateSession(ctx context.Context, user *model.User, param *parameters.CreateChatSession) (*model.ChatSession, error)
	SendMessage(ctx *gin.Context, user *model.User, param *parameters.CreateChatMessageRequest) (*responder.Manager, error)
	FeedbackMessage(ctx *gin.Context, user *model.User, id int64, vote bool) (*model.ChatMessageFeedback, error)
}
type chatBL struct {
	cfg                *Config
	chatRepo           repository.ChatRepository
	docRepo            repository.DocumentRepository
	personaRepo        repository.PersonaRepository
	embeddingModelRepo repository.EmbeddingModelRepository
	aiBuilder          *ai.Builder
	embedding          proto.EmbedServiceClient
	milvusClinet       storage.MilvusClient
}

func (b *chatBL) FeedbackMessage(ctx *gin.Context, user *model.User, id int64, vote bool) (*model.ChatMessageFeedback, error) {
	message, err := b.chatRepo.GetMessageByIDAndUserID(ctx, id, user.ID)
	if err != nil {
		return nil, err
	}
	feedback := message.Feedback
	if feedback == nil {
		feedback = &model.ChatMessageFeedback{
			ChatMessageID: message.ID,
			UserID:        user.ID,
		}
	}
	feedback.UpVotes = vote
	if err = b.chatRepo.MessageFeedback(ctx, feedback); err != nil {
		return nil, err
	}
	return feedback, nil
}

func (b *chatBL) SendMessage(ctx *gin.Context, user *model.User, param *parameters.CreateChatMessageRequest) (*responder.Manager, error) {
	chatSession, err := b.chatRepo.GetSessionByID(ctx.Request.Context(), user.ID, param.ChatSessionID.IntPart())
	if err != nil {
		return nil, err
	}
	em, err := b.embeddingModelRepo.GetDefault(ctx.Request.Context(), user.TenantID)
	if err != nil {
		zap.S().Errorf(err.Error())
		em = &model.EmbeddingModel{
			ModelID: b.cfg.DefaultEmbeddingModel,
		}
	}
	message := model.ChatMessage{
		ChatSessionID: chatSession.ID,
		Message:       param.Message,
		MessageType:   model.MessageTypeUser,
		TimeSent:      time.Now().UTC(),
	}
	noLLM := chatSession.Persona == nil
	if err = b.chatRepo.SendMessage(ctx.Request.Context(), &message); err != nil {
		return nil, err
	}
	aiClient := b.aiBuilder.New(chatSession.Persona.LLM)
	resp := responder.NewManager(
		responder.NewAIResponder(aiClient, b.chatRepo,
			b.embedding, b.milvusClinet, b.docRepo, em.ModelID),
	)

	go resp.Send(ctx, user, noLLM, &message, chatSession.Persona)
	return resp, nil
}

func (b *chatBL) GetSessions(ctx context.Context, user *model.User) ([]*model.ChatSession, error) {
	return b.chatRepo.GetSessions(ctx, user.ID)
}

func (b *chatBL) GetSessionByID(ctx context.Context, user *model.User, id int64) (*model.ChatSession, error) {
	result, err := b.chatRepo.GetSessionByID(ctx, user.ID, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *chatBL) CreateSession(ctx context.Context, user *model.User, param *parameters.CreateChatSession) (*model.ChatSession, error) {
	exists, err := b.personaRepo.IsExists(ctx, param.PersonaID.IntPart(), user.TenantID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, utils.ErrorBadRequest.New("persona is not exists")
	}
	session := model.ChatSession{
		UserID:       user.ID,
		Description:  param.Description,
		CreationDate: time.Now().UTC(),
		PersonaID:    param.PersonaID,
		OneShot:      param.OneShot,
	}
	if err = b.chatRepo.CreateSession(ctx, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func NewChatBL(
	cfg *Config,
	chatRepo repository.ChatRepository,
	personaRepo repository.PersonaRepository,
	docRepo repository.DocumentRepository,
	embeddingModelRepo repository.EmbeddingModelRepository,
	aiBuilder *ai.Builder,
	embedding proto.EmbedServiceClient,
	milvusClinet storage.MilvusClient,
) ChatBL {
	return &chatBL{
		cfg:                cfg,
		chatRepo:           chatRepo,
		personaRepo:        personaRepo,
		docRepo:            docRepo,
		embeddingModelRepo: embeddingModelRepo,
		aiBuilder:          aiBuilder,
		embedding:          embedding,
		milvusClinet:       milvusClinet,
	}
}
