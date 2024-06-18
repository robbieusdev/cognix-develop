package bll

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"time"
)

const minRefreshFreq = 3600

type (
	ConnectorBL interface {
		GetAll(ctx context.Context, user *model.User) ([]*model.Connector, error)
		GetByID(ctx context.Context, user *model.User, id int64) (*model.Connector, error)
		Create(ctx context.Context, user *model.User, param *parameters.CreateConnectorParam) (*model.Connector, error)
		Update(ctx context.Context, id int64, user *model.User, param *parameters.UpdateConnectorParam) (*model.Connector, error)
		Archive(ctx context.Context, user *model.User, id int64, restore bool) (*model.Connector, error)
	}
	connectorBL struct {
		connectorRepo repository.ConnectorRepository
		//messenger      messaging.Client
	}
)

func (c *connectorBL) Archive(ctx context.Context, user *model.User, id int64, restore bool) (*model.Connector, error) {
	connector, err := c.connectorRepo.GetByIDAndUser(ctx, user.TenantID, user.ID, id)
	if err != nil {
		return nil, err
	}
	if !(connector.UserID == user.ID || user.HasRoles(model.RoleAdmin, model.RoleSuperAdmin)) {
		return nil, utils.ErrorPermission.New("permission denied")
	}
	if !restore {
		connector.DeletedDate = pg.NullTime{time.Now().UTC()}
	} else {
		connector.DeletedDate = pg.NullTime{}
	}
	connector.LastUpdate = pg.NullTime{time.Now().UTC()}
	if err = c.connectorRepo.Update(ctx, connector); err != nil {
		return nil, err
	}
	return connector, nil
}

func NewConnectorBL(connectorRepo repository.ConnectorRepository) ConnectorBL {
	return &connectorBL{connectorRepo: connectorRepo}
}

func (c *connectorBL) Create(ctx context.Context, user *model.User, param *parameters.CreateConnectorParam) (*model.Connector, error) {

	tenantID := uuid.NullUUID{Valid: false}
	if param.Shared {
		tenantID.Valid = true
		tenantID.UUID = user.TenantID
	}
	conn := model.Connector{
		Name:                    param.Name,
		Type:                    model.SourceType(param.Source),
		ConnectorSpecificConfig: param.ConnectorSpecificConfig,
		RefreshFreq:             param.RefreshFreq,
		UserID:                  user.ID,
		TenantID:                tenantID,
		Status:                  model.ConnectorStatusReadyToProcessed,
		CreationDate:            time.Now().UTC(),
	}

	if conn.RefreshFreq < minRefreshFreq {
		conn.RefreshFreq = minRefreshFreq
	}
	if err := c.connectorRepo.Create(ctx, &conn); err != nil {
		return nil, err
	}
	return &conn, nil
}

func (c *connectorBL) Update(ctx context.Context, id int64, user *model.User, param *parameters.UpdateConnectorParam) (*model.Connector, error) {
	conn, err := c.connectorRepo.GetByIDAndUser(ctx, user.TenantID, user.ID, id)
	if err != nil {
		return nil, err
	}
	conn.ConnectorSpecificConfig = param.ConnectorSpecificConfig
	conn.Name = param.Name
	conn.RefreshFreq = param.RefreshFreq
	tenantID := uuid.NullUUID{Valid: false}
	if param.Shared {
		tenantID.Valid = true
		tenantID.UUID = user.TenantID
	}

	conn.TenantID = tenantID
	conn.LastUpdate = pg.NullTime{time.Now().UTC()}
	if param.Status != "" {
		conn.Status = param.Status
	}
	if conn.RefreshFreq < minRefreshFreq {
		conn.RefreshFreq = minRefreshFreq
	}
	if err = c.connectorRepo.Update(ctx, conn); err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *connectorBL) GetAll(ctx context.Context, user *model.User) ([]*model.Connector, error) {
	return c.connectorRepo.GetAllByUser(ctx, user.TenantID, user.ID)
}

func (c *connectorBL) GetByID(ctx context.Context, user *model.User, id int64) (*model.Connector, error) {
	return c.connectorRepo.GetByIDAndUser(ctx, user.TenantID, user.ID, id)
}
