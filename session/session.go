// Package provides additional layer between the client and message delivery to client instance

package session

type (
	Session interface {
		ClientId() string
		Send(action string, message interface{}) error
		SendErr(err error) error
	}
)
