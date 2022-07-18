package handlers

import (
	"github.com/Montesk/proofofwork/book"
	"github.com/Montesk/proofofwork/core/logger"
	"github.com/Montesk/proofofwork/pow/pow"
	powservice "github.com/Montesk/proofofwork/pow/service"
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
		log  logger.Logger
	}
)

func (h *handlers) All() map[Controller]Handler {
	return h.list
}

func New(log logger.Logger) Handlers {
	h := &handlers{
		pow:  powservice.New(),
		book: book.NewMemoryBook(),
		list: map[Controller]Handler{},
		log:  log,
	}

	h.setupRoutes()

	return h
}

func (h *handlers) setupRoutes() {
	h.list[ChallengeController] = BuildRoute[any](h.ChallengeHandler)
	h.list[ProveController] = BuildRoute[protocol.ProveController](h.ProveHandler)
}
