package server

import (
	"fmt"
	"github.com/Montesk/proofofwork/client"
	"github.com/Montesk/proofofwork/config"
	"github.com/Montesk/proofofwork/core/errors"
	"github.com/Montesk/proofofwork/core/logger"
	"github.com/Montesk/proofofwork/handlers"
	"github.com/Montesk/proofofwork/protocol"
	"github.com/Montesk/proofofwork/router"
	"github.com/Montesk/proofofwork/sessioner"
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
		config    config.Config
		router    router.Router
		sessioner sessioner.Sessioner
		listener  net.Listener
		shutdown  chan struct{}
		log       logger.Logger
	}
)

func New(cfg config.Config, rt router.Router, ss sessioner.Sessioner, shutdown chan struct{}, log logger.Logger) *server {
	return &server{
		sessioner: ss,
		config:    cfg,
		router:    rt,
		shutdown:  shutdown,
		log:       log,
	}
}

func (s *server) Run() error {
	listener, err := net.Listen(s.config.Protocol(), s.config.Port())
	if err != nil {
		return err
	}

	s.listener = listener

	s.registerHandlers()

	return nil
}

func (s *server) Close() error {
	if !s.ready() {
		return ErrServerNotReady
	}

	s.log.Infof("server %s closed", s.config.Port())

	return s.listener.Close()
}

func (s *server) Listen() error {
	if !s.ready() {
		return ErrServerNotReady
	}

	go s.listen()

	<-s.shutdown

	s.log.Info("received application shutdown signal")
	s.log.Info("tcp server closed for new connections")

	return nil
}

func (s *server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.log.Errorf("listen connection err %v", err)
			return
		}

		requestsChan := make(chan protocol.ClientMessage)
		errorsChan := make(chan error)

		cl := client.New[protocol.ClientMessage](conn, s.config.ReadTimeout(), requestsChan, errorsChan, s.log)

		err = s.registerClient(cl, conn)
		if err != nil {
			s.log.Errorf("can't register client %s err %v", cl.Addr(), err)
			continue
		}

		_, err = conn.Write([]byte(fmt.Sprintf("client %s connected \n", cl.Addr())))
		if err != nil {
			s.log.Errorf("send client %s ping message err %v", cl.Addr(), err)
		}

		s.log.Debugf("client: %s connected", cl.Addr())

		go s.handle(cl, requestsChan, errorsChan)
	}
}

func (s *server) ready() bool {
	return s.listener != nil
}

func (s *server) registerClient(client client.Client, conn net.Conn) error {
	return s.sessioner.Register(client, conn)
}

func (s *server) handle(cl client.Client, requests chan protocol.ClientMessage, errors chan error) {
	disconnect := make(chan struct{})

	defer func() {
		_ = cl.Close() // :WARNING: connection might be already closed so that's why we are in this defer; skip error already closed err logging
		disconnect <- struct{}{}
	}()

	go s.handleRequests(cl, requests, errors, disconnect)

	err := cl.Listen()
	if err != nil {
		s.log.Errorf("client %s listen err conn close %v", cl.Addr(), err)
	}
}

func (s *server) handleRequests(cl client.Client, requests chan protocol.ClientMessage, errors chan error, disconnect chan struct{}) {
	for {
		select {
		case req := <-requests:
			ses, err := s.sessioner.Session(cl.Addr())
			if err != nil {
				s.log.Errorf("error retrieve client session: %v client %v", err, cl.Addr())
				continue
			}

			err = s.router.Handle(req.Controller, ses, req.Message)
			if err != nil {
				s.log.Errorf("handler message err: %v client %v", err, cl.Addr())
			}
		case errMsg := <-errors:
			ses, err := s.sessioner.Session(cl.Addr())
			if err != nil {
				s.log.Errorf("sending error response to client %s err %v", cl.Addr(), err)
				continue
			}

			err = ses.SendErr(errMsg)
			if err != nil {
				s.log.Errorf("send err to client %s err %v", cl.Addr(), err)
			}
		case <-disconnect:
			err := s.sessioner.Unregister(cl)
			if err != nil {
				s.log.Errorf("client unregister err: %v client %v", err, cl.Addr())
			}

			return
		}
	}
}

func (s *server) registerHandlers() {
	h := handlers.New(s.log)

	for controller, handler := range h.All() {
		s.router.Register(string(controller), handler)
	}
}
