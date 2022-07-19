// Package provides clients storage

package sessioner

import (
	"github.com/Montesk/proofofwork/client"
	"github.com/Montesk/proofofwork/session"
	"net"
)

type (
	Sessioner interface {
		Session(clientId string) (session.Session, error)
		Register(client client.Client, conn net.Conn) error
		Unregister(client client.Client) error
		Close()
	}
)
