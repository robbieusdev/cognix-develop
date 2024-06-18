package repository

import (
	"cognix.ch/api/v2/core/utils"
	"go.uber.org/fx"
)

var DatabaseModule = fx.Options(
	fx.Provide(func() (*Config, error) {
		cfg := Config{}
		err := utils.ReadConfig(&cfg)
		return &cfg, err
	},
		NewDatabase,
	),
)

var RepositoriesModule = fx.Options(
	fx.Provide(
		NewUserRepository,
		NewConnectorRepository,
		NewLLMRepository,
		NewPersonaRepository,
		NewChatRepository,
		NewDocumentRepository,
		NewEmbeddingModelRepository,
		NewTenantRepository,
	),
)
