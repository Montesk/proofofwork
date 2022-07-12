package config

type (
	config struct {
		port     string
		protocol string
	}

	Config interface {
		Port() string
		Protocol() string
	}
)

func (c config) Port() string {
	return c.port
}

func (c config) Protocol() string {
	return c.protocol
}

func New(protocol, port string) Config {
	return &config{
		protocol: protocol,
		port:     port,
	}
}
