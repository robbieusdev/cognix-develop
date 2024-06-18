package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"context"
	"github.com/google/uuid"
)

type (
	Youtube struct {
		Base
		param *YoutubeParameters
	}
	YoutubeParameters struct {
		URL string `url:"url"`
	}
)

func (c *Youtube) PrepareTask(ctx context.Context, task Task) error {
	sessionID := uuid.NullUUID{
		UUID:  uuid.New(),
		Valid: true,
	}
	if len(c.model.Docs) == 0 {
		doc, ok := c.model.DocsMap[c.param.URL]
		if !ok {
			doc = &model.Document{
				SourceID:        c.param.URL,
				ConnectorID:     c.Base.model.ID,
				URL:             c.param.URL,
				Signature:       "",
				ChunkingSession: sessionID,
			}
			c.model.Docs = append(c.model.Docs, doc)
		}
	}
	// ignore  file that was analyzed
	if c.model.Status == model.ConnectorStatusError || c.model.Status == model.ConnectorStatusSuccess {
		return nil
	}
	return task.RunSemantic(ctx, &proto.SemanticData{
		Url:            c.param.URL,
		ConnectorId:    c.model.ID.IntPart(),
		FileType:       proto.FileType_YT,
		CollectionName: c.model.CollectionName(),
		ModelName:      c.model.User.EmbeddingModel.ModelID,
		ModelDimension: int32(c.model.User.EmbeddingModel.ModelDim),
	})
}

func (c *Youtube) Execute(ctx context.Context, param map[string]string) chan *Response {
	go func() {
		close(c.resultCh)
	}()
	return c.resultCh
}

func NewYoutube(connector *model.Connector) (Connector, error) {
	youtube := Youtube{}
	youtube.Base.Config(connector)
	youtube.param = &YoutubeParameters{}
	if err := connector.ConnectorSpecificConfig.ToStruct(youtube.param); err != nil {
		return nil, err
	}

	return &youtube, nil
}
