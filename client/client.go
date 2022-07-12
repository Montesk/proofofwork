package client

import (
	"bufio"
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

	client struct {
		conn        net.Conn
		readTimeout time.Duration
	}
)

func New(conn net.Conn, timeout time.Duration) Client {
	return &client{
		conn:        conn,
		readTimeout: timeout,
	}
}

func (c *client) Listen() error {
	err := c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	if err != nil {
		return err
	}

	for {
		reader := bufio.NewReader(c.conn)

		msg, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		log.Printf("recieved message from client msg %s", msg)
	}
}

func (c *client) Close() error {
	log.Printf("client connection closed %s", c.Addr())
	return c.conn.Close()
}

func (c *client) Addr() string {
	return c.conn.RemoteAddr().String()
}
