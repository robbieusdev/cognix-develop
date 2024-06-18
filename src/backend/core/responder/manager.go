package responder

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"sync"
)

type Manager struct {
	ch         chan *Response
	wg         *sync.WaitGroup
	responders []ChatResponder
}

func (m *Manager) Send(cx context.Context,
	user *model.User,
	noLLM bool,
	parentMessage *model.ChatMessage,
	persona *model.Persona) {
	for _, responder := range m.responders {
		m.wg.Add(1)
		go responder.Send(cx, m.ch, m.wg, user, noLLM, parentMessage, persona)
	}
	m.wg.Wait()
	close(m.ch)
}

func (m *Manager) Receive() (*Response, bool) {
	for response := range m.ch {
		return response, true
	}
	return &Response{
		IsValid: true,
		Type:    ResponseEnd,
	}, false
}

func NewManager(responders ...ChatResponder) *Manager {
	return &Manager{
		ch:         make(chan *Response, 1),
		wg:         &sync.WaitGroup{},
		responders: append([]ChatResponder{}, responders...),
	}
}
