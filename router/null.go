package router

import (
	"encoding/json"
	"github.com/Montesk/proofofwork/session"
)

type (
	null struct{}
)

func NewNull() *null {
	return &null{}
}

func (r *null) Register(path string, handler Route)                                    {}
func (r *null) Handle(path string, ses session.Session, message json.RawMessage) error { return nil }
