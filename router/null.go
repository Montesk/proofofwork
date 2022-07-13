package router

type (
	null struct{}
)

func Null() Router {
	return &router{
		routes: map[string]func(){},
	}
}

func (r *null) Register(path string, handler func()) {}
func (r *null) Handle(path string) error             { return nil }
