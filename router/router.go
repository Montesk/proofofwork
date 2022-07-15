package router

import (
	"encoding/json"
	"github.com/Montesk/proofofwork/errors"
	"github.com/Montesk/proofofwork/handlers"
	"github.com/Montesk/proofofwork/session"
	"log"
)

const (
	ErrRouteNotRegistered = errors.String("unknown route")
)

type (
	Router interface {
		Register(path string, handler handlers.Handler)
		Handle(path string, ses session.Session, message json.RawMessage) error
	}

	router struct {
		routes map[string]handlers.Handler
	}
)

func New() Router {
	return &router{
		routes: map[string]handlers.Handler{},
	}
}

func (r *router) Register(path string, handler handlers.Handler) {
	_, found := r.routes[path]
	if found {
		log.Fatalf("route %s already registered", path)
	}

	r.routes[path] = handler
}

func (r *router) Handle(path string, ses session.Session, message json.RawMessage) error {
	handler, ok := r.routes[path]
	if !ok {
		log.Printf("route %s not registered", path)
		return ErrRouteNotRegistered
	}

	err := handler.Prepare(message)
	if err != nil {
		return err
	}

	handler.Call(ses)

	return nil
}
