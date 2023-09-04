package common

import (
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
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
	Path          string
	Filename      string
	ChunkSize     int
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
func (c *Client) StartClientLoop() {
	// autoincremental msgID to identify every message sent
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	filename := strings.Replace(c.config.Filename, "{id}", c.config.ID, 1)

	chunkReader, err := newChunkReader(filename, c.config.ChunkSize)
	if err != nil {
		log.Errorf("action: new_chunk_reader | result: fail | error: %v", err)
		return
	}

	c.createClientSocket()

	for {
		select {
		case <-sigs:
			log.Infof("action: signal_received")
			break
		default:
		}

		lines, err := chunkReader.read()
		if err != nil {
			log.Errorf("action: chunk_reader_read | result: fail | error: %v", err)
			break
		}

		lines_with_agency := bets_array_add_agency(c.config.ID, lines)
		bets := bets_from_array(lines_with_agency)
		chunk := bets_to_chunk(bets)

		err = sendString(c.conn, chunk)
		if err != nil {
			log.Errorf("action: send_string | result: fail | error: %v", err)
			break
		}

		if len(lines) == 0 {
			resp, err := recvString(c.conn)
			if err != nil {
				log.Errorf("action: recv_string | result: fail | error: %v", err)
			}
			log.Infof("action: recv_string | result: success | msg: %v", resp)
			break
		}
		// Wait a time between sending one message and the next one
		//time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)

	id, err := strconv.Atoi(c.config.ID)
	if err != nil {
		log.Errorf("action: atoi | result: fail | error: %v", err)
	}

	err = sendU32(c.conn, uint32(id))
	if err != nil {
		log.Errorf("action: send_u32 | result: fail | error: %v", err)
	}

	chunk, err := recvString(c.conn)
	if err != nil {
		log.Errorf("action: recv string chunk | result: fail | error: %v", err)
	}

	bets := bet_documents_from_chunk(chunk)
	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", len(bets))

	log.Infof("action: lottery_finished | result: success | client_id: %v", c.config.ID)
}
