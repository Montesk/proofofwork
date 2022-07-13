package router

import (
	"github.com/faraway/wordofwisdom/errors"
	"log"
)

const (
	ErrRouteNotRegistered = errors.String("unknown route")
)

type (
	Router interface {
		Register(path string, handler func())
		Handle(path string) error
	}

	router struct {
		routes map[string]func()
	}
)

func New(routes map[string]func()) Router {
	return &router{
		routes: routes,
	}
}

func (r *router) Register(path string, handler func()) {
	r.routes[path] = handler
}

func (r *router) Handle(path string) error {
	handler, ok := r.routes[path]
	if !ok {
		log.Printf("route %s not registered", path)
		return ErrRouteNotRegistered
	}

	handler()

	return nil
}
