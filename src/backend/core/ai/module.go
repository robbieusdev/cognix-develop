package ai

import (
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/utils"
	"fmt"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ChunkingModule = fx.Options(
	fx.Provide(func() (*ChunkingConfig, error) {
		cfg := ChunkingConfig{}
		if err := utils.ReadConfig(&cfg); err != nil {
			return nil, err
		}
		if err := cfg.Validate(); err != nil {
			return nil, err
		}
		return &cfg, nil
	},
		newChunking,
	),
)

func newChunking(cfg *ChunkingConfig) Chunking {
	if cfg.Strategy == StrategyLLM {
		return NewLLMChunking()
	}
	return NewStaticChunking(cfg)
}

var EmbeddingModule = fx.Options(
	fx.Provide(func() (*EmbeddingConfig, error) {
		cfg := EmbeddingConfig{}
		if err := utils.ReadConfig(&cfg); err != nil {
			return nil, err
		}
		return &cfg, nil
	},
		newEmbeddingGRPCClient),
)

func newEmbeddingGRPCClient(cfg *EmbeddingConfig) (proto.EmbedServiceClient, error) {
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.EmbedderHost, cfg.EmbedderPort), dialOptions...)
	if err != nil {
		return nil, err
	}
	return proto.NewEmbedServiceClient(conn), nil
}
