package config

import "time"

type (
	mock struct {
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

func NewMockConfig(protocol, port string, readTimeout time.Duration, clientsNumber int) *mock {
	return &mock{
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

func (m *mock) Protocol() string {
	return m.server.protocol
}

func (m *mock) Port() string {
	return m.server.port
}

func (m *mock) ReadTimeout() time.Duration {
	return m.server.readTimeout
}

func (m *mock) POWClients() (number int) {
	return m.clients.number
}
