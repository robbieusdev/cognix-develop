package bll

import (
	"cognix.ch/api/v2/core/utils"
	"go.uber.org/fx"
)

type Config struct {
	RedirectURL                string `env:"REDIRECT_URL"`
	DefaultEmbeddingModel      string `env:"DEFAULT_EMBEDDING_MODEL" envDefault:"paraphrase-multilingual-mpnet-base-v2"`
	DefaultEmbeddingVectorSize int    `env:"DEFAULT_EMBEDDING_VECTOR_SIZE" envDefault:"768"`
}

var BLLModule = fx.Options(
	fx.Provide(func() (*Config, error) {
		cfg := Config{}
		err := utils.ReadConfig(&cfg)
		return &cfg, err
	}),
	fx.Provide(
		NewConnectorBL,
		NewAuthBL,
		NewPersonaBL,
		NewChatBL,
		NewDocumentBL,
		NewEmbeddingModelBL,
		NewTenantBL,
	),
)
