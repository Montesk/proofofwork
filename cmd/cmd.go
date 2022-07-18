// Package that provides implementation(s) of command line argument read on app start

package cmd

import (
	"github.com/Montesk/proofofwork/core/logger"
	"time"
)

type (
	Cmd interface {
		Port() int
		Protocol() string
		ReadTimeout() time.Duration
		POWClients() (number int)
		LogLevel() logger.Level
	}
)
