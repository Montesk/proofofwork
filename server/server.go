package server

import (
	"fmt"
	"git.ll-games.com/backend/daily/go1.17rc1/src/log"
	"github.com/faraway/wordofwisdom/client"
	"github.com/faraway/wordofwisdom/config"
	"github.com/faraway/wordofwisdom/errors"
	"net"
	"time"
)

const (
	DefaultReadTimeout = 10 * time.Second

	ErrServerNotReady = errors.String("server not ready")
)

type (
	Server interface {
		Run() error
		Listen() error
		Close() error
	}

	server struct {
		cfg      config.Config
		listener net.Listener
	}
)

func New(cfg config.Config) *server {
	return &server{
		cfg: cfg,
	}
}

func (s *server) Run() error {
	listener, err := net.Listen(s.cfg.Protocol(), s.cfg.Port())
	if err != nil {
		return err
	}

	log.Print("server started port: ", s.cfg.Port())

	s.listener = listener

	return nil
}

func (s *server) Close() error {
	if !s.ready() {
		return ErrServerNotReady
	}

	log.Printf("server %s closed", s.cfg.Port())

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

		cl := client.New(conn, s.cfg.ReadTimeout())

		log.Printf("client: %s connected", cl.Addr())

		go func() {
			defer cl.Close()
			cl.Listen()
		}()
	}
}

func (s *server) ready() bool {
	return s.listener != nil
}
