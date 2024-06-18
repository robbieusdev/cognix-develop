package main

import (
	"cognix.ch/api/v2/core/ai"
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/storage"
	"cognix.ch/api/v2/core/utils"
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Config struct {
	OAuthURL         string `env:"OAUTH_URL,required"`
	SubscriptionName string `env:"SUBSCRIPTION_NAME,required"`
}

var Module = fx.Options(
	repository.DatabaseModule,
	messaging.NatsModule,
	storage.MinioModule,
	storage.MilvusModule,
	ai.ChunkingModule,
	fx.Provide(
		func() (*Config, error) {
			cfg := Config{}
			err := utils.ReadConfig(&cfg)
			if err != nil {
				zap.S().Errorf(err.Error())
				return nil, err
			}
			return &cfg, nil
		},
		repository.NewConnectorRepository,
		repository.NewDocumentRepository,
		repository.NewEmbeddingModelRepository,
		NewExecutor,
	),
	fx.Invoke(RunServer),
)

func RunServer(lc fx.Lifecycle, executor *Executor) error {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go executor.run(executor.msgClient.StreamConfig().ConnectorStreamName, executor.msgClient.StreamConfig().ConnectorStreamSubject, executor.runConnector)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			executor.msgClient.Close()
			return nil
		},
	})
	return nil
}
