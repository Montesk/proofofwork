package router

import (
	"encoding/json"
	"github.com/Montesk/proofofwork/core/errors"
	"github.com/Montesk/proofofwork/session"
)

const (
	ErrRouteNotRegistered = errors.String("unknown route")
)

type (
	Router interface {
		Register(path string, handler Route)
		Handle(path string, ses session.Session, message json.RawMessage) error
	}

	Route interface {
		Call(ses session.Session)
		Prepare(message []byte) error
	}
)
