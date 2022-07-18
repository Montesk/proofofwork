// Package provide Router implementation of ControllerName -> HandlerCb() dictionary

package router

import (
	"encoding/json"
	"github.com/Montesk/proofofwork/core/logger"
	"github.com/Montesk/proofofwork/session"
)

type (
	controllers struct {
		log    logger.Logger
		routes map[string]Route
	}
)

func NewControllers(log logger.Logger) *controllers {
	return &controllers{
		routes: map[string]Route{},
		log:    log,
	}
}

func (r *controllers) Register(path string, handler Route) {
	_, found := r.routes[path]
	if found {
		r.log.Fatalf("route %s already registered", path)
	}

	r.routes[path] = handler
}

func (r *controllers) Handle(path string, ses session.Session, message json.RawMessage) error {
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
