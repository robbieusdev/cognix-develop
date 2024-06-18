package main

import (
	"cognix.ch/api/v2/core/connector"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/repository"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

const (
	ConnectorSchedulerSpan = "connector-scheduler"
)

type (
	trigger struct {
		messenger      messaging.Client
		connectorRepo  repository.ConnectorRepository
		docRepo        repository.DocumentRepository
		tracer         trace.Tracer
		connectorModel *model.Connector
		fileSizeLimit  int
	}
)

func (t *trigger) Do(ctx context.Context) error {
	// if connector is new or
	// todo we need figure out how to use multiple  orchestrators instances
	// one approach could be that this method will extract top x rows from the database
	// and it will book them

	if t.connectorModel.User == nil || t.connectorModel.User.EmbeddingModel == nil {
		return fmt.Errorf("embedding model is not configured for %s", t.connectorModel.Name)
	}
	zap.S().Debugf("\n------------  %s\nlast %v refresh Freq %d \nnext %v\nnow  %v\nlast+refreshFreq > now %v\n------------- ",
		t.connectorModel.Name,
		t.connectorModel.LastUpdate.UTC(),
		t.connectorModel.RefreshFreq,
		t.connectorModel.LastUpdate.UTC().Add(time.Duration(t.connectorModel.RefreshFreq)*time.Second),
		time.Now().UTC(),
		t.connectorModel.LastUpdate.UTC().Add(time.Duration(t.connectorModel.RefreshFreq)*time.Second).Before(time.Now().UTC()))

	if t.connectorModel.LastUpdate.IsZero() ||
		t.connectorModel.LastUpdate.UTC().Add(time.Duration(t.connectorModel.RefreshFreq)*time.Second).Before(time.Now().UTC()) {
		ctx, span := t.tracer.Start(ctx, ConnectorSchedulerSpan)
		defer span.End()
		span.SetAttributes(attribute.Int64(model.SpanAttributeConnectorID, t.connectorModel.ID.IntPart()))
		span.SetAttributes(attribute.String(model.SpanAttributeConnectorSource, string(t.connectorModel.Type)))

		//if err := t.updateStatus(ctx, model.ConnectorStatusPending); err != nil {
		//	span.RecordError(err)
		//	return err
		//}
		connWF, err := connector.New(t.connectorModel)
		if err != nil {
			return err
		}

		if err = connWF.PrepareTask(ctx, t); err != nil {
			span.RecordError(err)
			zap.S().Errorf("failed to prepare task for connector %s[%d]: %v", t.connectorModel.Name, t.connectorModel.ID.IntPart(), err)
			if errr := t.updateStatus(ctx, model.ConnectorStatusUnableProcess); errr != nil {
				span.RecordError(errr)
			}
			return err
		}
	}
	return nil
}

// RunSemantic send message to semantic service
func (t *trigger) RunSemantic(ctx context.Context, data *proto.SemanticData) error {

	if t.connectorModel.Type == model.SourceTypeWEB ||
		t.connectorModel.Type == model.SourceTypeYoutube ||
		t.connectorModel.Type == model.SourceTypeFile {
		doc := t.connectorModel.Docs[0]
		var err error
		// create or update document in database
		if doc.ID.IntPart() != 0 {
			err = t.docRepo.Update(ctx, doc)
		} else {
			err = t.docRepo.Create(ctx, doc)
		}
		if err != nil {
			zap.S().Errorf("update document failed %v", err)
			return err
		}
		data.DocumentId = doc.ID.IntPart()
	}
	if err := t.updateStatus(ctx, model.ConnectorStatusPending); err != nil {
		return err
	}
	zap.S().Infof("send message to semantic %s", t.connectorModel.Name)
	buf, _ := json.Marshal(data)
	zap.S().Debugf(" message payload %s", string(buf))
	return t.messenger.Publish(ctx, t.messenger.StreamConfig().SemanticStreamName,
		t.messenger.StreamConfig().SemanticStreamSubject, data)
}

// RunConnector send message to connector service
func (t *trigger) RunConnector(ctx context.Context, data *proto.ConnectorRequest) error {
	data.Params[connector.ParamFileLimit] = fmt.Sprintf("%d", t.fileSizeLimit)
	if err := t.updateStatus(ctx, model.ConnectorStatusPending); err != nil {
		return err
	}
	zap.S().Infof("send message to connector %s", t.connectorModel.Name)
	return t.messenger.Publish(ctx, t.messenger.StreamConfig().ConnectorStreamName,
		t.messenger.StreamConfig().ConnectorStreamSubject, data)
}

func (t *trigger) UpToDate(ctx context.Context) error {
	// may be to be implemented in future
	return nil
}

func NewTrigger(messenger messaging.Client,
	connectorRepo repository.ConnectorRepository,
	docRepo repository.DocumentRepository,
	connectorModel *model.Connector,
	fileSizeLimit int) *trigger {
	return &trigger{
		messenger:      messenger,
		connectorRepo:  connectorRepo,
		docRepo:        docRepo,
		connectorModel: connectorModel,
		fileSizeLimit:  fileSizeLimit,
		tracer:         otel.Tracer(model.TracerConnector),
	}
}

// update status of connector in database
func (t *trigger) updateStatus(ctx context.Context, status string) error {
	t.connectorModel.Status = status
	t.connectorModel.LastUpdate = pg.NullTime{time.Now().UTC()}
	return t.connectorRepo.Update(ctx, t.connectorModel)
}
