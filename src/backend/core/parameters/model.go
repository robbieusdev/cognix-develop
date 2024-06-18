package parameters

import (
	"cognix.ch/api/v2/core/model"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type LoginParam struct {
	RedirectURL string `form:"redirect_url"`
}

type OAuthParam struct {
	Action   string `json:"action,omitempty"`
	TenantID string `json:"tenant_id,omitempty"`
	Role     string `json:"role,omitempty"`
	Email    string `json:"email,omitempty"`
}

type InviteParam struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (v InviteParam) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Email, validation.Required, is.Email),
		validation.Field(&v.Role, validation.Required, validation.In(model.RoleSuperAdmin, model.RoleAdmin, model.RoleUser)),
	)
}

type ArchivedParam struct {
	Archived bool `form:"archived"`
}

type CreateConnectorParam struct {
	Name                    string        `json:"name,omitempty"`
	Source                  string        `json:"source,omitempty"`
	ConnectorSpecificConfig model.JSONMap `json:"connector_specific_config,omitempty"`
	RefreshFreq             int           `json:"refresh_freq,omitempty"`
	Shared                  bool          `json:"shared,omitempty"`
	Disabled                bool          `json:"disabled,omitempty"`
}

func (v CreateConnectorParam) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Source, validation.Required,
			validation.By(func(value interface{}) error {
				if st, ok := model.AllSourceTypes[model.SourceType(v.Source)]; !ok || !st.IsImplemented {
					return fmt.Errorf("invalid source type")
				}
				return nil
			})))
}

type UpdateConnectorParam struct {
	Name                    string        `json:"name,omitempty"`
	ConnectorSpecificConfig model.JSONMap `json:"connector_specific_config,omitempty"`
	RefreshFreq             int           `json:"refresh_freq,omitempty"`
	Shared                  bool          `json:"shared,omitempty"`
	Status                  string        `json:"status"`
}

func (v UpdateConnectorParam) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Name, validation.Required),
		validation.Field(&v.ConnectorSpecificConfig, validation.Required),
		validation.Field(&v.RefreshFreq, validation.Required),
		validation.Field(&v.Status, validation.In("", model.ConnectorStatusReadyToProcessed)),
	)
}

type AddUserParam struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (v AddUserParam) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Email, validation.Required, is.Email),
		validation.Field(&v.Role, validation.Required, validation.In(model.RoleSuperAdmin, model.RoleUser, model.RoleAdmin)),
	)
}

type EditUserParam struct {
	Role string `json:"role"`
}

func (v EditUserParam) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Role, validation.Required, validation.In(model.RoleSuperAdmin, model.RoleUser, model.RoleAdmin)),
	)
}
