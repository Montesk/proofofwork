package session

import (
	"encoding/json"
	"github.com/Montesk/proofofwork/protocol"
	"net"
)

type (
	networked struct {
		clientId string
		conn     net.Conn
	}
)

func (n *networked) ClientId() string {
	return n.clientId
}

func (n *networked) Send(action string, message interface{}) error {
	response := protocol.Action{
		Action:  action,
		Message: message,
	}

	raw, err := json.Marshal(response)
	if err != nil {
		return err
	}

	_, err = n.conn.Write(append(raw, '\n'))
	if err != nil {
		return err
	}

	return nil
}

func (n *networked) SendErr(err error) error {
	response := protocol.Error{
		Action: "error",
		Error:  err.Error(),
	}

	raw, err := json.Marshal(response)
	if err != nil {
		return err
	}

	_, err = n.conn.Write(append(raw, '\n'))
	if err != nil {
		return err
	}

	return nil
}

func NewNetworked(clientId string, conn net.Conn) *networked {
	return &networked{
		clientId: clientId,
		conn:     conn,
	}
}
