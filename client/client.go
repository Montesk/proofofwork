package client

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
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
	}
)

func New[T any](conn net.Conn, timeout time.Duration, requests chan T) Client {
	return &client[T]{
		conn:        conn,
		readTimeout: timeout,
		requests:    requests,
	}
}

func (c *client[T]) Listen() error {
	err := c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	if err != nil {
		return err
	}

	for {
		reader := bufio.NewReader(c.conn)

		raw, err := reader.ReadBytes('\n')
		if clientDisconnected(err) {
			return nil
		} else if err != nil {
			log.Printf("error parse message from client %v", err)
			continue
		}

		err = c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
		if err != nil {
			return err
		}

		msg := *new(T)

		err = json.Unmarshal(raw, &msg)
		if err != nil {
			log.Printf("error parse message from client %v", err)
			continue
		}

		log.Printf("recieved message from client msg %v", msg)

		c.requests <- msg
	}
}

func (c *client[T]) Close() error {
	log.Printf("client connection closed %s", c.Addr())
	return c.conn.Close()
}

func (c *client[T]) Addr() string {
	return c.conn.RemoteAddr().String()
}

func clientDisconnected(got error) bool {
	var err net.Error
	return errors.Is(got, io.EOF) || (errors.As(got, &err) && err.Timeout())
}
