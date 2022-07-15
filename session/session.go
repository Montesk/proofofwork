package session

import (
	"encoding/json"
	"github.com/Montesk/wordofwisdom/protocol"
	"net"
)

type (
	Session interface {
		ClientId() string
		Send(action string, message interface{}) error
	}

	session struct {
		clientId string
		conn     net.Conn
	}
)

func (s *session) ClientId() string {
	return s.clientId
}

func (s *session) Send(action string, message interface{}) error {
	response := protocol.Action{
		Action:  action,
		Message: message,
	}

	raw, err := json.Marshal(response)
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
