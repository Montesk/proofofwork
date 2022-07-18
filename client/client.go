package client

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/Montesk/proofofwork/core/logger"
	"io"
	"net"
	"time"
)

type (
	Client interface {
		Listen() error
		Close() error
		Addr() string
	}

	// typed parameter T must be supported to json.Unmarshal value e.g. interface, map[string]interface or concrete struct
	client[T any] struct {
		conn        net.Conn
		readTimeout time.Duration
		requests    chan T
		errors      chan error
		log         logger.Logger
	}
)

func New[T any](conn net.Conn, timeout time.Duration, requests chan T, errors chan error, log logger.Logger) Client {
	return &client[T]{
		conn:        conn,
		readTimeout: timeout,
		requests:    requests,
		errors:      errors,
		log:         log,
	}
}

func (c *client[T]) Listen() error {
	err := c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	if err != nil {
		c.log.Errorf("error set deadline for client %v err %v", c.Addr(), err)
		return err
	}

	reader := bufio.NewReader(c.conn)

	for {
		raw, err := reader.ReadBytes('\n')
		if clientDisconnected(err) {
			return nil
		} else if err != nil {
			c.log.Errorf("error parse message from client %v err %v", c.Addr(), err)
			continue
		}

		err = c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
		if err != nil {
			c.log.Errorf("error set deadline for client %v err %v", c.Addr(), err)
			return err
		}

		msg := *new(T)

		err = json.Unmarshal(raw, &msg)
		if err != nil {
			c.log.Errorf("parse message from client err %v", err)
			c.errors <- err
			continue
		}

		c.log.Debugf("received message from client msg %s", c.Addr())

		c.requests <- msg
	}
}

func (c *client[T]) Close() error {
	c.log.Debugf("client connection closed %s", c.Addr())
	return c.conn.Close()
}

func (c *client[T]) Addr() string {
	return c.conn.RemoteAddr().String()
}

func clientDisconnected(got error) bool {
	var err net.Error
	return errors.Is(got, io.EOF) || (errors.As(got, &err) && err.Timeout())
}
