package router

import (
	"encoding/json"
	"github.com/Montesk/proofofwork/core/errors"
	"github.com/Montesk/proofofwork/core/logger"
	"github.com/Montesk/proofofwork/handlers"
	"github.com/Montesk/proofofwork/session"
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
		log    logger.Logger
		routes map[string]handlers.Handler
	}
)

func New(log logger.Logger) Router {
	return &router{
		routes: map[string]handlers.Handler{},
		log:    log,
	}
}

func (r *router) Register(path string, handler handlers.Handler) {
	_, found := r.routes[path]
	if found {
		r.log.Fatalf("route %s already registered", path)
	}

	r.routes[path] = handler
}

func (r *router) Handle(path string, ses session.Session, message json.RawMessage) error {
	handler, ok := r.routes[path]
	if !ok {
		r.log.Errorf("route %s not registered", path)
		return ErrRouteNotRegistered
	}

	err := handler.Prepare(message)
	if err != nil {
		return err
	}

	handler.Call(ses)

	return nil
}
