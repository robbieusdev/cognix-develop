package parameters

import validation "github.com/go-ozzo/ozzo-validation/v4"

type EmbeddingModelParam struct {
	ModelID   string `json:"model_id,omitempty"`
	ModelName string `json:"model_name,omitempty"`
	ModelDim  int    `json:"model_dim,omitempty"`
	URL       string `json:"url,omitempty"`
	IsActive  bool   `json:"is_active,omitempty"`
}

func (v EmbeddingModelParam) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.ModelID, validation.Required))
}
