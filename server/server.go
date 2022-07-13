package server

import (
	"fmt"
	"git.ll-games.com/backend/daily/go1.17rc1/src/log"
	"github.com/faraway/wordofwisdom/client"
	"github.com/faraway/wordofwisdom/config"
	"github.com/faraway/wordofwisdom/errors"
	"github.com/faraway/wordofwisdom/protocol"
	"github.com/faraway/wordofwisdom/router"
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
		config   config.Config
		router   router.Router
		listener net.Listener
	}
)

func New(config config.Config, router router.Router) *server {
	return &server{
		config: config,
		router: router,
	}
}

func (s *server) Run() error {
	listener, err := net.Listen(s.config.Protocol(), s.config.Port())
	if err != nil {
		return err
	}

	log.Print("server started port: ", s.config.Port())

	s.listener = listener

	return nil
}

func (s *server) Close() error {
	if !s.ready() {
		return ErrServerNotReady
	}

	log.Printf("server %s closed", s.config.Port())

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

		requests := make(chan protocol.Client)

		cl := client.New[protocol.Client](conn, s.config.ReadTimeout(), requests)

		log.Printf("client: %s connected", cl.Addr())

		go s.handle(cl, requests)
	}
}

func (s *server) ready() bool {
	return s.listener != nil
}

func (s *server) handle(cl client.Client, requests chan protocol.Client) {
	disconnect := make(chan struct{})

	defer func() {
		cl.Close()
		disconnect <- struct{}{}
	}()

	go s.handleRequests(cl, requests, disconnect)

	cl.Listen()
}

func (s *server) handleRequests(cl client.Client, requests chan protocol.Client, disconnect chan struct{}) {
	for {
		select {
		case req := <-requests:
			err := s.router.Handle(req.Controller)
			if err != nil {
				log.Printf("error handler message err: %v client %v", err, cl.Addr())
			}
		case <-disconnect:
			return
		}
	}
}
