// Package is the implementation of client with some work to suggest nonce parts

package client

import (
	"fmt"
	"github.com/Montesk/proofofwork/core/errors"
	"github.com/Montesk/proofofwork/core/slice"
	"github.com/Montesk/proofofwork/pow/pow"
	"strconv"
)

const MaxTries = pow.NonceMax + 1

const (
	ErrTooManyTries = errors.String("out of try limits")
)

type (
	Client interface {
		Suggest(challenge string, tries int) (int, error)
	}

	client struct {
		id         string
		service    pow.POW
		nonceTries []int
	}
)

func New(clientId string, pow pow.POW) Client {
	return &client{
		id:         clientId,
		service:    pow,
		nonceTries: []int{},
	}
}

func (c *client) Suggest(challenge string, tries int) (int, error) {
	if tries > MaxTries {
		return 0, ErrTooManyTries
	}

	success := c.service.Prove(c.id, c.work(challenge, tries))

	if !success {
		tries += 1

		return c.Suggest(challenge, tries)
	}

	return tries, nil
}

func (c *client) work(challenge string, try int) (suggest string) {
	if slice.Contains[int](try, c.nonceTries) {
		return c.work(challenge, try)
	}

	c.nonceTries = append(c.nonceTries, try)

	return fmt.Sprintf("%s%s%s", challenge, pow.HashDelimiter, pow.Encode([]string{strconv.Itoa(try)}))
}
