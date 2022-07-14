package session

import (
	"encoding/json"
	"net"
)

type (
	Session interface {
		ClientId() string
		Send(message interface{}) error
	}

	session struct {
		clientId string
		conn     net.Conn
	}
)

func (s *session) ClientId() string {
	return s.clientId
}

func (s *session) Send(message interface{}) error {
	raw, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = s.conn.Write(raw)
	if err != nil {
		return err
	}

	return nil
}

func New(clientId string, conn net.Conn) Session {
	return &session{
		clientId: clientId,
		conn:     conn,
	}
}
