package logger

import "github.com/sirupsen/logrus"

type (
	debug struct {
		*logrus.Logger
	}
)

func NewDebug() *debug {
	l := logrus.New()
	l.SetLevel(DebugLevel)

	return &debug{
		l,
	}
}
