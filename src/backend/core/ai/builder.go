package ai

import (
	"cognix.ch/api/v2/core/model"
	"sync"
)

type Builder struct {
	clients map[int64]OpenAIClient
	mx      sync.Mutex
}

func NewBuilder() *Builder {
	return &Builder{clients: make(map[int64]OpenAIClient)}
}

func (b *Builder) New(llm *model.LLM) OpenAIClient {
	b.mx.Lock()
	defer b.mx.Unlock()
	if client, ok := b.clients[llm.ID.IntPart()]; ok {
		return client
	}
	client := NewOpenAIClient(llm.ModelID, llm.ApiKey)
	b.clients[llm.ID.IntPart()] = client
	return client
}
func (b *Builder) Invalidate(llm *model.LLM) {
	b.mx.Lock()
	delete(b.clients, llm.ID.IntPart())
	b.mx.Unlock()
}
