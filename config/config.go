// Package that provides implementation(s) application config build

package config

import (
	"github.com/Montesk/proofofwork/core/logger"
	"time"
)

const (
	DefaultPort             = 8001
	DefaultProtocol         = "tcp"
	DefaultLogLevel         = logger.DebugLevel
	DefaultReadTimeout      = 600 // seconds
	DefaultPOWClientsNumber = 100
)

type (
	Config interface {
		Protocol() string
		Port() string
		ReadTimeout() time.Duration
		POWClients() (number int)
		LogLevel() logger.Level
	}
)
