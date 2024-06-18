package ai

import (
	_ "github.com/deluan/flowllm/llms/openai"
)

type (
	EmbeddingConfig struct {
		EmbedderHost string `env:"EMBEDDER_GRPC_HOST,required"`
		EmbedderPort int    `env:"EMBEDDER_GRPC_PORT,required"`
	}
)
