package server

import (
	"fmt"
	"git.ll-games.com/backend/daily/go1.17rc1/src/log"
	"github.com/Montesk/proofofwork/client"
	"github.com/Montesk/proofofwork/config"
	"github.com/Montesk/proofofwork/errors"
	"github.com/Montesk/proofofwork/handlers"
	"github.com/Montesk/proofofwork/protocol"
	"github.com/Montesk/proofofwork/router"
	"github.com/Montesk/proofofwork/sessioner"
	"net"
	"time"
)

const (
	DefaultReadTimeout = 10 * time.Minute

	ErrServerNotReady = errors.String("server not ready")
)

type (
	Server interface {
		Run() error
		Listen() error
		Close() error
	}

	server struct {
		config    config.Config
		router    router.Router
		sessioner sessioner.Sessioner
		listener  net.Listener
	}
)

func New(cfg config.Config, rt router.Router, ss sessioner.Sessioner) *server {
	return &server{
		sessioner: ss,
		config:    cfg,
		router:    rt,
	}
}

func (s *server) Run() error {
	listener, err := net.Listen(s.config.Protocol(), s.config.Port())
	if err != nil {
		return err
	}

	log.Print("server started port: ", s.config.Port())

	s.listener = listener

	s.registerHandlers()

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

		requests := make(chan protocol.ClientMessage)

		cl := client.New[protocol.ClientMessage](conn, s.config.ReadTimeout(), requests)

		err = s.registerClient(cl, conn)
		if err != nil {
			log.Printf("can't register client %s err %v", cl.Addr(), err)
			continue
		}

		log.Printf("client: %s connected", cl.Addr())

		go s.handle(cl, requests)
	}
}

func (s *server) ready() bool {
	return s.listener != nil
}

func (s *server) registerClient(client client.Client, conn net.Conn) error {
	return s.sessioner.Register(client, conn)
}

func (s *server) handle(cl client.Client, requests chan protocol.ClientMessage) {
	disconnect := make(chan struct{})

	defer func() {
		cl.Close()
		disconnect <- struct{}{}
	}()

	go s.handleRequests(cl, requests, disconnect)

	cl.Listen()
}

func (s *server) handleRequests(cl client.Client, requests chan protocol.ClientMessage, disconnect chan struct{}) {
	for {
		select {
		case req := <-requests:
			ses, err := s.sessioner.Session(cl.Addr())
			if err != nil {
				log.Printf("error retrieve client session: %v client %v", err, cl.Addr())
				continue
			}

			err = s.router.Handle(req.Controller, ses, req.Message)
			if err != nil {
				log.Printf("error handler message err: %v client %v", err, cl.Addr())
			}
		case <-disconnect:
			err := s.sessioner.Unregister(cl)
			if err != nil {
				log.Printf("error client unregister err: %v client %v", err, cl.Addr())
			}

			return
		}
	}
}

func (s *server) registerHandlers() {
	h := handlers.New()

	for controller, handler := range h.All() {
		s.router.Register(string(controller), handler)
	}
}
