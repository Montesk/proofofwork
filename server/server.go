package server

import (
	"fmt"
	"git.ll-games.com/backend/daily/go1.17rc1/src/log"
	"github.com/faraway/wordofwisdom/config"
	"github.com/faraway/wordofwisdom/errors"
	"net"
)

const (
	ErrServerNotReady = errors.String("server not ready")
)

type (
	Server interface {
		Run() error
		Listen() error
		Close() error
	}

	server struct {
		protocol string
		port     string
		listener net.Listener
	}
)

func New(cfg config.Config) *server {
	return &server{
		protocol: cfg.Protocol(),
		port:     cfg.Port(),
	}
}

func (s *server) Run() error {
	listener, err := net.Listen(s.protocol, s.port)
	if err != nil {
		return err
	}

	s.listener = listener

	return nil
}

func (s *server) Close() error {
	if !s.ready() {
		return ErrServerNotReady
	}

	return s.listener.Close()
}

func (s *server) Listen() error {
	if !s.ready() {
		return ErrServerNotReady
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return fmt.Errorf("accept listener err %v", err)
		}

		log.Printf("received connection %v", conn)

		conn.Close()

		log.Printf("connection closed")
	}
}

func (s *server) ready() bool {
	return s.listener != nil
}
