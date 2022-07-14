package handlers

import "github.com/faraway/wordofwisdom/session"

type (
	Controller string
)

type (
	Handlers interface {
		All() map[Controller]handler
	}

	handler func(ses session.Session)

	handlers struct {
		list map[Controller]handler
	}
)

func (h *handlers) All() map[Controller]handler {
	return h.list
}

func New() Handlers {
	return &handlers{
		list: map[Controller]handler{
			ChallengeController: ChallengeHandler,
		},
	}
}
