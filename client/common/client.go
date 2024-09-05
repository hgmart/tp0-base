package common

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID             string
	ServerAddress  string
	LoopAmount     int
	LoopPeriod     time.Duration
	PayloadMaxSize int
	ZipBatchPath   string
}

// Client Entity that encapsulates how
type Client struct {
	config    ClientConfig
	conn      net.Conn
	isEnabled bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config:    config,
		isEnabled: true,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	if !c.isEnabled {
		return nil
	}

	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)

		c.isEnabled = false
	}
	c.conn = conn
	return nil
}

func (c *Client) closeSocketConnection() {
	if c.conn == nil {
		return
	}

	err := c.conn.Close()
	if err != nil {
		log.Errorf("action: closing_connection | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
}

func (c *Client) waitServerReponse(signals chan os.Signal) {
	if !c.isEnabled {
		return
	}

	ch_msg := make(chan string)
	ch_err := make(chan error)

	go func() {
		log.Infof("action: waiting_message | result: success | client_id: %v | msg: waiting echo message",
			c.config.ID,
		)

		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			ch_err <- err
			return
		}

		ch_msg <- msg
	}()

	select {
	case msg := <-ch_msg:
		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
			c.config.ID,
			msg,
		)

	case err := <-ch_err:
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)

	case received_signal := <-signals:
		log.Infof("action: sigterm_signal | result: success | client_id: %v | msg: %v signal received",
			c.config.ID,
			received_signal,
		)

		c.isEnabled = false
	}
}

func (c *Client) delayMessageDelivery(signals chan os.Signal) {
	if !c.isEnabled {
		return
	}

	select {
	// Waits a time between sending one message and the next one
	case currentTime := <-time.After(c.config.LoopPeriod):
		log.Infof("action: sleep_time | result: success | client_id: %v | msg: awaiked at %v",
			c.config.ID,
			currentTime,
		)

	case received_signal := <-signals:
		log.Debugf("action: sigterm_signal | result: success | client_id: %v | msg: %v signal received",
			c.config.ID,
			received_signal,
		)

		c.isEnabled = false
	}
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)

	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; c.isEnabled && msgID <= c.config.LoopAmount; msgID++ {
		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		// TODO: Modify the send to avoid short-write
		fmt.Fprintf(
			c.conn,
			"[CLIENT %v] Message N°%v\n",
			c.config.ID,
			msgID,
		)

		c.waitServerReponse(signals)
		c.closeSocketConnection()
		c.delayMessageDelivery(signals)
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
