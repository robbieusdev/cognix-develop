package repository

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/lib/pq"
	"time"
)

type (
	DocumentRepository interface {
		FindByConnectorIDAndUser(ctx context.Context, user *model.User, connectorID int64) ([]*model.Document, error)
		FindByConnectorID(ctx context.Context, connectorID int64) ([]*model.Document, error)
		FindByID(ctx context.Context, id int64) (*model.Document, error)
		Create(ctx context.Context, document *model.Document) error
		Update(ctx context.Context, document *model.Document) error
		DeleteByIDS(ctx context.Context, ids ...int64) error
	}
	documentRepository struct {
		db *pg.DB
	}
)

func (r *documentRepository) FindByID(ctx context.Context, id int64) (*model.Document, error) {
	var doc model.Document
	if err := r.db.WithContext(ctx).Model(&doc).Where("id = ?", id).Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "document not found")
	}
	return &doc, nil
}

func (r *documentRepository) FindByConnectorID(ctx context.Context, connectorID int64) ([]*model.Document, error) {
	documents := make([]*model.Document, 0)
	if err := r.db.WithContext(ctx).Model(&documents).
		Join("INNER JOIN connectors c ON c.id = connector_id").
		Where("connector_id = ?", connectorID).
		Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find documents ")
	}
	return documents, nil
}

func (r *documentRepository) FindByConnectorIDAndUser(ctx context.Context, user *model.User, connectorID int64) ([]*model.Document, error) {
	documents := make([]*model.Document, 0)
	if err := r.db.WithContext(ctx).Model(&documents).
		Join("INNER JOIN connectors c ON c.id = connector_id").
		Where("connector_id = ?", connectorID).
		Where("c.tenant_id = ?", user.TenantID).
		WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			return query.WhereOr("c.user_id = ? ", user.ID).
				WhereOr("c.shared = ?", true), nil
		}).Select(); err != nil {
		return nil, utils.NotFound.Wrap(err, "can not find documents ")
	}
	return documents, nil
}

func (r *documentRepository) Create(ctx context.Context, document *model.Document) error {
	document.CreationDate = time.Now().UTC()
	stm := r.db.WithContext(ctx).Model(document)
	if !document.ParentID.Valid {
		stm = stm.ExcludeColumn("parent_id")
	}

	if _, err := stm.Insert(); err != nil {
		return utils.Internal.Wrapf(err, "can not insert document [%s]", err.Error())
	}
	return nil
}

func (r *documentRepository) Update(ctx context.Context, document *model.Document) error {
	stm := r.db.WithContext(ctx).Model(document).Where("id = ? ", document.ID)
	if !document.ParentID.Valid {
		stm = stm.ExcludeColumn("parent_id")
	}

	document.LastUpdate = pg.NullTime{time.Now().UTC()}
	if _, err := stm.Update(); err != nil {
		return utils.Internal.Wrapf(err, "can not update document [%s]", err.Error())
	}
	return nil
}

func (r *documentRepository) DeleteByIDS(ctx context.Context, ids ...int64) error {
	if _, err := r.db.WithContext(ctx).Model(&model.Document{}).
		Where("id = any ?", pq.Array(ids)).Delete(); err != nil {
		return utils.Internal.Wrapf(err, "can not delete documents [%s]", err.Error())
	}
	return nil
}

func NewDocumentRepository(db *pg.DB) DocumentRepository {
	return &documentRepository{db: db}
}
