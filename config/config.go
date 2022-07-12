package config

import "time"

type (
	Config interface {
		Protocol() string
		Port() string
		ReadTimeout() time.Duration
	}

	config struct {
		protocol    string
		port        string
		readTimeout time.Duration
	}
)

func New(protocol, port string, readTimeout time.Duration) Config {
	return &config{
		protocol:    protocol,
		port:        port,
		readTimeout: readTimeout,
	}
}

func (c config) Protocol() string {
	return c.protocol
}

func (c config) Port() string {
	return c.port
}

func (c config) ReadTimeout() time.Duration {
	return c.readTimeout
}
