package sessioner

import (
	"github.com/Montesk/proofofwork/client"
	"github.com/Montesk/proofofwork/core/errors"
	"github.com/Montesk/proofofwork/session"
	"net"
	"sync"
)

const (
	ErrClientNotFound          = errors.String("client not found")
	ErrClientAlreadyRegistered = errors.String("client already registered in the system")
)

type (
	Sessioner interface {
		Session(clientId string) (session.Session, error)
		Register(client client.Client, conn net.Conn) error
		Unregister(client client.Client) error
	}

	sessioner struct {
		mu      *sync.Mutex
		storage map[string]struct {
			client  client.Client
			session session.Session
		}
	}
)

func New() Sessioner {
	return &sessioner{
		mu: new(sync.Mutex),
		storage: map[string]struct {
			client  client.Client
			session session.Session
		}{},
	}
}

func (s *sessioner) Session(clientId string) (session.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cl, found := s.storage[clientId]
	if !found {
		return nil, ErrClientNotFound
	}

	return cl.session, nil
}

func (s *sessioner) Register(cl client.Client, conn net.Conn) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, found := s.storage[cl.Addr()]
	if found {
		return ErrClientAlreadyRegistered
	}

	s.storage[cl.Addr()] = struct {
		client  client.Client
		session session.Session
	}{
		client:  cl,
		session: session.New(cl.Addr(), conn),
	}

	return nil
}

func (s *sessioner) Unregister(client client.Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, found := s.storage[client.Addr()]
	if !found {
		return ErrClientNotFound
	}

	delete(s.storage, client.Addr())

	return nil
}
