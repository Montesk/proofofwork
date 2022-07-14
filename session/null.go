package session

type null struct{}

func (s *null) ClientId() string                              { return "" }
func (s *null) Send(action string, message interface{}) error { return nil }

func Null() Session {
	return &null{}
}
