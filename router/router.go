package router

import (
	"github.com/faraway/wordofwisdom/errors"
	"github.com/faraway/wordofwisdom/session"
	"log"
)

const (
	ErrRouteNotRegistered = errors.String("unknown route")
)

type (
	Router interface {
		Register(path string, handler func(ses session.Session))
		Handle(path string, ses session.Session) error
	}

	router struct {
		routes map[string]func(ses session.Session)
	}
)

func New() Router {
	return &router{
		routes: map[string]func(ses session.Session){},
	}
}

func (r *router) Register(path string, handler func(ses session.Session)) {
	r.routes[path] = handler
}

func (r *router) Handle(path string, ses session.Session) error {
	handler, ok := r.routes[path]
	if !ok {
		log.Printf("route %s not registered", path)
		return ErrRouteNotRegistered
	}

	handler(ses)

	return nil
}
