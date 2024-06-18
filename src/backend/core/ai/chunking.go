package ai

import (
	"context"
)

type staticChunking struct {
	chunkSize int
	overlap   int
}

func (s staticChunking) Split(ctx context.Context, text string) ([]string, error) {
	var chunks []string
	for {
		if len(text) < s.chunkSize {
			chunks = append(chunks, text)
			break
		}
		chunks = append(chunks, text[:s.chunkSize])
		text = text[s.chunkSize-s.overlap:]
	}
	return chunks, nil
}

func NewStaticChunking(cfg *ChunkingConfig) Chunking {
	return &staticChunking{
		chunkSize: cfg.StaticChunkSize,
		overlap:   cfg.StaticChunkOverlap,
	}
}

func NewLLMChunking() Chunking {
	return &staticChunking{
		chunkSize: MaxChunkSize,
		overlap:   20,
	}
}
