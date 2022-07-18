package logger

import "github.com/sirupsen/logrus"

type (
	base struct {
		*logrus.Logger
	}
)

func NewBase(level Level) *base {
	l := logrus.New()
	l.SetLevel(level)

	return &base{l}
}
