package router

import "github.com/faraway/wordofwisdom/session"

type (
	null struct{}
)

func Null() Router {
	return &router{
		routes: map[string]func(ses session.Session){},
	}
}

func (r *null) Register(path string, handler func()) {}
func (r *null) Handle(path string) error             { return nil }
