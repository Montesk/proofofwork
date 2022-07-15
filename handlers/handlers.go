package handlers

import (
	"github.com/Montesk/proofofwork/book"
	"github.com/Montesk/proofofwork/pow"
	"github.com/Montesk/proofofwork/protocol"
)

type (
	Controller string
)

type (
	Handlers interface {
		All() map[Controller]Handler
	}

	handlers struct {
		pow  pow.POW
		book book.Book
		list map[Controller]Handler
	}
)

func (h *handlers) All() map[Controller]Handler {
	return h.list
}

func New() Handlers {
	h := &handlers{
		pow:  pow.New(),
		book: book.New(),
		list: map[Controller]Handler{},
	}

	h.setupRoutes()

	return h
}

func (h *handlers) setupRoutes() {
	h.list[ChallengeController] = BuildRoute[any](h.ChallengeHandler)
	h.list[ProveController] = BuildRoute[protocol.ProveController](h.ProveHandler)
}
