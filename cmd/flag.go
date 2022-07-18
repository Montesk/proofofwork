package cmd

import (
	"flag"
	"github.com/Montesk/proofofwork/config"
	"sync"
	"time"
)

type (
	flagCmd struct {
		port             int
		protocol         string
		readTimeout      time.Duration
		powClientsNumber int
	}
)

var once = new(sync.Once)

func NewFlagCmd() *flagCmd {
	var port, timeout, clients int
	var protocol string

	// command line args must get default values only once
	once.Do(func() {
		flag.IntVar(&port, "port", config.DefaultPort, "run the application server on specified port")
		flag.IntVar(&timeout, "timeout", config.DefaultReadTimeout, "specify the application server read timeout in seconds")
		flag.StringVar(&protocol, "protocol", config.DefaultProtocol, "run the application server with specified protocol")
		flag.IntVar(&clients, "clients", config.DefaultPOWClientsNumber, "specify the number of clients to start proof of work challenge")

		flag.Parse()
	})

	return &flagCmd{
		port:             port,
		protocol:         protocol,
		powClientsNumber: clients,
		readTimeout:      time.Duration(timeout) * time.Second,
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
