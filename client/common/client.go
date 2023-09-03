package common

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Fatalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(bet *Bet) {
	// autoincremental msgID to identify every message sent
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

loop:
	// Send messages if the loopLapse threshold has not been surpassed
	for timeout := time.After(c.config.LoopLapse); ; {
		select {
		case <-timeout:
			log.Infof("action: timeout_detected | result: success | client_id: %v",
				c.config.ID,
			)
			break loop
		case <-sigs:
			log.Infof("action: signal_received")
			break loop
		default:
		}

		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		if err := sendBet(c.conn, bet); err != nil {
			log.Errorf("action: send_bet | result: fail | document: %v | number: %v",
				bet.document,
				bet.number,
			)
			return
		}

		if _, err := recvString(c.conn); err != nil {
			log.Errorf("action: recv_string | result: fail | document: %v | number: %v",
				bet.document,
				bet.number,
			)
			return
		}

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
