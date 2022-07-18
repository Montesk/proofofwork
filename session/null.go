package session

type null struct{}

func (s *null) ClientId() string                              { return "" }
func (s *null) Send(action string, message interface{}) error { return nil }
func (s *null) SendErr(err error) error                       { return nil }

func Null() *null {
	return &null{}
}
