package pow

import (
	"fmt"
	"github.com/faraway/wordofwisdom/errors"
	"strconv"
)

// mock implementation of client with some work to suggest nonce part of challenge in SHA-1

const MaxTries = nonceMax + 1

const (
	ErrTooManyTries = errors.String("out of try limits")
)

type (
	Client interface {
		Work(challenge string, try int) (suggest string)
		Suggest(challenge string, tries int) (int, error)
	}

	client struct {
		id         string
		service    POW
		nonceTries []int
	}
)

func (c *client) Work(challenge string, try int) (suggest string) {
	if intSliceContains(try, c.nonceTries) {
		return c.Work(challenge, try)
	}

	c.nonceTries = append(c.nonceTries, try)

	return fmt.Sprintf("%s%s%s", challenge, hashDelimiter, prepareHash([]string{strconv.Itoa(try)}))
}

func (c *client) Suggest(challenge string, tries int) (int, error) {
	if tries > MaxTries {
		return 0, ErrTooManyTries
	}

	success := c.service.Prove(c.id, c.Work(challenge, tries))

	if !success {
		tries += 1

		return c.Suggest(challenge, tries)
	}

	return tries, nil
}

func newMockClient(clientId string, pow POW) Client {
	return &client{
		id:         clientId,
		service:    pow,
		nonceTries: []int{},
	}
}
