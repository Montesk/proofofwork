package errors

type String string

func (s String) Error() string {
	return string(s)
}
