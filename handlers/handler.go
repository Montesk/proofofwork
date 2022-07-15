package handlers

import (
	"encoding/json"
	"github.com/faraway/wordofwisdom/session"
)

type (
	Handler interface {
		Call(ses session.Session)
		Prepare(message []byte) error
	}

	handler[T any] struct {
		cb    func(ses session.Session, param T)
		param T
	}
)

func BuildRoute[T any](cb func(ses session.Session, param T)) Handler {
	return &handler[T]{cb: cb}
}

func (h *handler[T]) Call(ses session.Session) {
	h.cb(ses, h.param)
}

func (h *handler[T]) Prepare(message []byte) error {
	if message == nil {
		return nil
	}

	err := json.Unmarshal(message, &h.param)
	if err != nil {
		return err
	}

	return nil
}
