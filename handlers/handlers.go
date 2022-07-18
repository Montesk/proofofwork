package handlers

import (
	"github.com/Montesk/proofofwork/book"
	"github.com/Montesk/proofofwork/core/logger"
	"github.com/Montesk/proofofwork/pow/pow"
	powservice "github.com/Montesk/proofofwork/pow/service"
	"github.com/Montesk/proofofwork/protocol"
	"github.com/Montesk/proofofwork/router"
)

type (
	Controller string
)

type (
	Handlers interface {
		All() map[Controller]router.Route
	}

	handlers struct {
		pow  pow.POW
		book book.Book
		list map[Controller]router.Route
		log  logger.Logger
	}
)

func (h *handlers) All() map[Controller]router.Route {
	return h.list
}

func New(log logger.Logger) Handlers {
	h := &handlers{
		pow:  powservice.New(),
		book: book.NewMemoryBook(),
		list: map[Controller]router.Route{},
		log:  log,
	}

	h.setupRoutes()

	return h
}

func (h *handlers) setupRoutes() {
	h.list[ChallengeController] = router.BuildRoute[any](h.ChallengeHandler)
	h.list[ProveController] = router.BuildRoute[protocol.ProveController](h.ProveHandler)
}
