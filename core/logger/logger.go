// Package provides constructor helpers for "logrus" logger initializing

package logger

import "github.com/sirupsen/logrus"

// see logrus.Level for each level detailed description
const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type (
	Logger logrus.FieldLogger

	Level = logrus.Level
)

func ParseLevel(lvl string) (Level, error) {
	return logrus.ParseLevel(lvl)
}
