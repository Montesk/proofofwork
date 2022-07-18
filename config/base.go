package config

import (
	"github.com/Montesk/proofofwork/core/logger"
	"time"
)

type (
	base struct {
		app struct {
			logLevel logger.Level
		}
		server struct {
			protocol    string
			port        string
			readTimeout time.Duration
		}
		clients struct {
			number int
		}
	}
)

func NewBase(protocol, port string, readTimeout time.Duration, clientsNumber int, logLevel logger.Level) *base {
	return &base{
		app: struct{ logLevel logger.Level }{logLevel: logLevel},
		server: struct {
			protocol    string
			port        string
			readTimeout time.Duration
		}{
			protocol:    protocol,
			port:        port,
			readTimeout: readTimeout,
		},

		clients: struct {
			number int
		}{
			number: clientsNumber,
		},
	}
}

func (b *base) Protocol() string {
	return b.server.protocol
}

func (b *base) Port() string {
	return b.server.port
}

func (b *base) ReadTimeout() time.Duration {
	return b.server.readTimeout
}

func (b *base) POWClients() (number int) {
	return b.clients.number
}

func (b *base) LogLevel() logger.Level {
	return b.app.logLevel
}
