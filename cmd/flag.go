package cmd

import (
	"flag"
	"github.com/Montesk/proofofwork/config"
	"github.com/Montesk/proofofwork/core/logger"
	"sync"
	"time"
)

type (
	flagCmd struct {
		port             int
		protocol         string
		readTimeout      time.Duration
		powClientsNumber int
		logLevel         logger.Level
	}
)

var initialized bool
var once = new(sync.Once)

func NewFlagCmd(log logger.Logger) *flagCmd {
	if initialized {
		log.Warn("command args[] has been already initialized")
	}

	var port, timeout, clients int
	var protocol, logLevelArg string

	// command line args must get default values only once
	once.Do(func() {
		initialized = true

		flag.IntVar(&port, "port", config.DefaultPort, "run the application server on specified port")
		flag.IntVar(&timeout, "timeout", config.DefaultReadTimeout, "specify the application server read timeout in seconds")
		flag.StringVar(&protocol, "protocol", config.DefaultProtocol, "run the application server with specified protocol")
		flag.StringVar(&logLevelArg, "log_level", config.DefaultLogLevel.String(), "run the application with specified log level")
		flag.IntVar(&clients, "clients", config.DefaultPOWClientsNumber, "specify the number of clients to start proof of work challenge")

		flag.Parse()
	})

	logLevel, err := logger.ParseLevel(logLevelArg)
	if err != nil {
		log.Errorf("failed to parse log level got arg %s", logLevel)
		log.Infof("setting default log level %v...", config.DefaultLogLevel.String())
		logLevel = config.DefaultLogLevel
	}

	return &flagCmd{
		port:             port,
		protocol:         protocol,
		powClientsNumber: clients,
		readTimeout:      time.Duration(timeout) * time.Second,
		logLevel:         logLevel,
	}
}

func (f *flagCmd) Port() int {
	return f.port
}

func (f *flagCmd) Protocol() string {
	return f.protocol
}

func (f *flagCmd) ReadTimeout() time.Duration {
	return f.readTimeout
}

func (f *flagCmd) POWClients() (number int) {
	return f.powClientsNumber
}

func (f *flagCmd) LogLevel() logger.Level {
	return f.logLevel
}
