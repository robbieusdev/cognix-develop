package responder

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"sync"
)

const (
	ResponseMessage  = "message"
	ResponseError    = "error"
	ResponseDocument = "document"
	ResponseEnd      = "end"
)

type Response struct {
	IsValid  bool
	Type     string
	Message  *model.ChatMessage
	Document *model.DocumentResponse
	Err      error
}

type ChatResponder interface {
	Send(cx context.Context, ch chan *Response, wg *sync.WaitGroup, user *model.User, noLLM bool, parentMessage *model.ChatMessage, persona *model.Persona)
}
