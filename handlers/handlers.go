package handlers

import (
	"github.com/faraway/wordofwisdom/book"
	"github.com/faraway/wordofwisdom/pow"
	"github.com/faraway/wordofwisdom/session"
)

type (
	Controller string
)

type (
	Handlers interface {
		All() map[Controller]handler
	}

	handler func(ses session.Session)

	handlers struct {
		pow  pow.POW
		book book.Book
		list map[Controller]handler
	}
)

func (h *handlers) All() map[Controller]handler {
	return h.list
}

func New() Handlers {
	h := &handlers{
		pow:  pow.New(),
		book: book.New(),
		list: map[Controller]handler{},
	}

	h.setupRoutes()

	return h
}

func (h *handlers) setupRoutes() {
	h.list[ChallengeController] = h.ChallengeHandler
}
