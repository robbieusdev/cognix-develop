package model

import "github.com/google/uuid"

type Tenant struct {
	tableName     struct{}  `pg:"tenants"`
	ID            uuid.UUID `json:"id`
	Name          string    `json:"name"`
	Configuration JSONMap   `json:"configuration" pg:"type:jsonb"`
}
