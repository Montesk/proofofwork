package router

import (
	"github.com/faraway/wordofwisdom/handlers"
)

type (
	null struct{}
)

func Null() Router {
	return &router{
		routes: map[string]handlers.Handler{},
	}
}

func (r *null) Register(path string, handler func()) {}
func (r *null) Handle(path string) error             { return nil }
