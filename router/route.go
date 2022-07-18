package router

import (
	"encoding/json"
	"github.com/Montesk/proofofwork/session"
)

type (
	route[T any] struct {
		cb    func(ses session.Session, param T)
		param T
	}
)

func BuildRoute[T any](cb func(ses session.Session, param T)) Route {
	return &route[T]{cb: cb}
}

func (h *route[T]) Call(ses session.Session) {
	h.cb(ses, h.param)
}

func (h *route[T]) Prepare(message []byte) error {
	if message == nil {
		return nil
	}

	err := json.Unmarshal(message, &h.param)
	if err != nil {
		return err
	}

	return nil
}
