// Implementation of client with TCP read/write generate and proof

package networked

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Montesk/proofofwork/config"
	"github.com/Montesk/proofofwork/core/errors"
	"github.com/Montesk/proofofwork/core/logger"
	"github.com/Montesk/proofofwork/pow/pow"
	"github.com/Montesk/proofofwork/protocol"
	"net"
	"strconv"
	"time"
)

const (
	ErrClientNotConnected = errors.String("client not connected to server")
)

type (
	networked struct {
		config config.Config
		conn   net.Conn
		log    logger.Logger
	}
)

func New(cfg config.Config, log logger.Logger) pow.POW {
	port, err := strconv.Atoi(cfg.Port())
	if err != nil {
		log.Fatal(err)
	}

	cl := &networked{
		config: cfg,
		log:    log,
	}

	conn, err := net.DialTCP(cfg.Protocol(), nil, &net.TCPAddr{
		Port: port,
	})
	if err != nil {
		log.Errorf("client failed to establish connection %v", err)
		return cl
	}

	cl.conn = conn

	go cl.listenServerStatus(cl.conn.LocalAddr().String())

	return cl
}

func (n *networked) Generate(clientId string) (string, error) {
	if n.conn == nil {
		n.log.Errorf("client %s err %v", clientId, ErrClientNotConnected)
		return "", ErrClientNotConnected
	}

	msg := protocol.ClientMessage{
		Controller: protocol.ChallengeController,
	}

	raw, err := json.Marshal(msg)
	if err != nil {
		n.log.Error(err)
		return "", err
	}

	_, err = n.conn.Write(append(raw, '\n'))
	if err != nil {
		n.log.Errorf("client %s write error %v", clientId, err)
		return "", err
	}

	result, err := waitForMessage[protocol.ChallengeActionMsg](n.conn)
	if err != nil {
		n.log.Errorf("client wait generate message %s error %v", clientId, err)
		return "", err
	}

	n.log.Debugf("client N %s received challenge", clientId)

	return result.Challenge, err
}

func (n *networked) Prove(clientId, hash string) (success bool) {
	if n.conn == nil {
		n.log.Errorf("client %s err %v", clientId, ErrClientNotConnected)
		return
	}

	msg := protocol.ClientMessage{
		Controller: protocol.ProveController,
		Message:    []byte(fmt.Sprintf(`{ "suggest": "%s" }`, hash)),
	}

	raw, err := json.Marshal(msg)
	if err != nil {
		n.log.Error(err)
		return false
	}

	_, err = n.conn.Write(append(raw, '\n'))
	if err != nil {
		n.log.Errorf("client %s write error %v", clientId, err)
		return false
	}

	result, err := waitForMessage[protocol.ProveActionMsg](n.conn)
	if err != nil {
		n.log.Errorf("client wait prove message %s error %v", clientId, err)
		return false
	}

	if result.Success {
		n.log.Infof("client %s N %s successfully decode message: %s", clientId, n.conn.LocalAddr(), result.Message)
		// :WARING: can't close connection here as new connection can take the same system port if system runs concurrently
		// :NOTE: connection in the end will be closed by the server
	}

	return result.Success
}

func (n *networked) listenServerStatus(clientId string) {
	reader := bufio.NewReader(n.conn)

	for {
		raw, err := reader.ReadBytes('\n')
		if err != nil {
			n.log.Errorf("client read bytes listening server status err %v", err)
			return
		}

		wrapper := protocol.Action{}
		err = json.Unmarshal(raw, &wrapper)
		if err != nil || wrapper.Action != protocol.InfoAction {
			// wait only for server info actions
			continue
		}

		res, ok := wrapper.Message.(map[string]any)
		if !ok {
			continue
		}

		rawFromMap, _ := json.Marshal(res)

		msg := protocol.ServerInfoActionMsg{}
		err = json.Unmarshal(rawFromMap, &msg)
		if err != nil {
			n.log.Errorf("server status unmarshal error", err)
			continue
		}

		n.log.Infof("client %s got server message: %v", clientId, msg.Info)
		return
	}
}

func waitForMessage[T any](conn net.Conn) (T, error) {
	reader := bufio.NewReader(conn)

	for {
		raw, err := reader.ReadBytes('\n')
		if err != nil {
			return *new(T), err
		}

		select {
		case <-time.After(10 * time.Second):
			return *new(T), fmt.Errorf("waiting message deadline")
		default:
			wrapper := protocol.Action{}

			err = json.Unmarshal(raw, &wrapper)
			if err != nil {
				// :KLUDGE: as we wait for specific message passed from typed parameter let's say it's ok to ignore unmarshal err for integration test purposes
				continue
			}

			expect := *new(T)

			res, ok := wrapper.Message.(map[string]any)
			if !ok {
				continue
			}

			rawFromMap, _ := json.Marshal(res)

			err = json.Unmarshal(rawFromMap, &expect)
			if err != nil {
				continue
			} else {
				return expect, nil
			}
		}

	}
}
