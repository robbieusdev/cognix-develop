package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type (
	File struct {
		Base
		param *FileParameters
		ctx   context.Context
	}
	FileParameters struct {
		FileName string `json:"file_name"`
		MIMEType string `json:"mime_type"`
	}
)

func (c *File) PrepareTask(ctx context.Context, task Task) error {

	link := fmt.Sprintf("minio:tenant-%s:%s", c.model.User.EmbeddingModel.TenantID.String(), c.param.FileName)
	if len(c.model.Docs) == 0 {
		c.model.Docs = append(c.model.Docs, &model.Document{
			SourceID:     link,
			ConnectorID:  c.model.ID,
			URL:          link,
			CreationDate: time.Now().UTC(),
			IsExists:     true,
		})
	}
	// ignore  file that was analyzed
	if c.model.Status == model.ConnectorStatusError || c.model.Status == model.ConnectorStatusSuccess {
		return nil
	}
	c.model.Docs[0].ChunkingSession = uuid.NullUUID{
		UUID:  uuid.New(),
		Valid: true,
	}
	return task.RunSemantic(ctx, &proto.SemanticData{
		Url:            link,
		DocumentId:     c.model.Docs[0].ID.IntPart(),
		ConnectorId:    c.model.ID.IntPart(),
		FileType:       0,
		CollectionName: c.model.CollectionName(),
		ModelName:      c.model.User.EmbeddingModel.ModelID,
		ModelDimension: int32(c.model.User.EmbeddingModel.ModelDim),
	})
}

func (c *File) Execute(ctx context.Context, param map[string]string) chan *Response {
	// do no used for this source
	c.ctx = ctx
	go func() {
		defer close(c.resultCh)
		if c.param == nil || c.param.FileName == "" {
			return
		}
		// check id document  already exists
		doc, ok := c.Base.model.DocsMap[c.param.FileName]
		url := fmt.Sprintf("minio:tenant-%s:%s", c.model.User.EmbeddingModel.TenantID, c.param.FileName)
		if !ok {
			doc = &model.Document{
				SourceID:    url,
				ConnectorID: c.model.ID,
				URL:         url,
				Signature:   "",
			}
			c.model.DocsMap[url] = doc
		}
		doc.IsExists = true
		if fileType, ok := supportedMimeTypes[c.param.MIMEType]; ok {
			c.resultCh <- &Response{
				URL:      url,
				SourceID: url,
				FileType: fileType,
			}
		} else {
			zap.S().Errorf("Upsupported file type : %s ", c.param.MIMEType)
		}

	}()
	return c.resultCh
}

// NewFile creates new instance of file connector.
func NewFile(connector *model.Connector) (Connector, error) {
	fileConn := File{}
	fileConn.Base.Config(connector)
	fileConn.param = &FileParameters{}
	if err := connector.ConnectorSpecificConfig.ToStruct(fileConn.param); err != nil {
		return nil, err
	}

	return &fileConn, nil
}
