package common

import (
	"io"
	"os"
	"os/signal"
	"syscall"
)

func (c *Client) receive_message(ch_msg chan []byte, ch_err chan error, chunk_number int) []byte {
	log.Infof("action: waiting_message | result: success | client_id: %v | chunk_number: %v | msg: waiting echo message",
		c.config.ID,
		chunk_number,
	)

	buffer := make([]byte, c.config.PayloadMaxSize)
	bytesReceived, err := c.conn.Read(buffer)

	log.Infof("action: bytes_received  | result: %v | client_id: %v | chunk_number: %v | bytes_number: %v | msg: %v",
		err == nil,
		c.config.ID,
		chunk_number,
		bytesReceived,
		buffer[:bytesReceived],
	)

	if err != nil && err != io.EOF {
		ch_err <- err
		return nil
	}

	dataReceived := buffer[:bytesReceived]

	if len(buffer) == bytesReceived {
		shortReadBytesReceived := c.receive_message(ch_msg, ch_err, chunk_number+1)
		dataReceived = append(buffer, shortReadBytesReceived...)
	}

	if chunk_number == 1 {
		ch_msg <- dataReceived
		return nil
	} else {
		return dataReceived
	}

}

func (c *Client) waitResponse(signals chan os.Signal) ([]byte, error) {
	if !c.isEnabled {
		return nil, nil
	}

	ch_msg := make(chan []byte)
	ch_err := make(chan error)

	go c.receive_message(ch_msg, ch_err, 1)

	select {
	case msg := <-ch_msg:
		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
			c.config.ID,
			msg,
		)
		return msg, nil

	case err := <-ch_err:
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return nil, err

	case received_signal := <-signals:
		log.Infof("action: sigterm_signal | result: success | client_id: %v | msg: %v signal received",
			c.config.ID,
			received_signal,
		)

		c.isEnabled = false
		return nil, nil
	}
}

func (c *Client) DefinePayloadSize(dataLength int) int {

	if c.config.PayloadMaxSize < dataLength {
		return c.config.PayloadMaxSize
	} else {
		return dataLength
	}

}

func (c *Client) SendData(data []byte, signals chan os.Signal, chunk_number int) (int, error) {

	if !c.isEnabled {
		return 0, nil
	}

	chunkSize := c.DefinePayloadSize(len(data))
	chunkData := data[:chunkSize]

	bytesWritten, err := c.conn.Write(chunkData)

	if err != nil {
		return -1, err
	}

	log.Infof("action: data_sent | client: %v | chunk_number: %v | chunk_size: %v | msg: %v", c.config.ID, chunk_number, len(chunkData), chunkData)

	if bytesWritten == len(data) {
		return bytesWritten, nil
	}

	bytesSent, err := c.SendData(data[bytesWritten:], signals, chunk_number+1)

	return bytesWritten + bytesSent, err
}

func (c *Client) SendBytes(chunkData []byte) (bool, []byte, error) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)

	go func() {
		received_signal := <-signals
		c.isEnabled = false

		log.Infof("action: sigterm_signal | result: success | client_id: %v | msg: %v signal received",
			c.config.ID,
			received_signal,
		)
	}()

	c.createClientSocket()
	log.Infof("action: message_to_be_sent | client: %v | chunk_size: %v | msg: %v", c.config.ID, len(chunkData), chunkData)

	bytesWritten, sentErr := c.SendData(chunkData, signals, 1)
	log.Infof("action: mensaje_enviado | result: %v | client: %v | chunk_size: %v | bytes_sent: %v", sentErr == nil, c.config.ID, len(chunkData), bytesWritten)

	if sentErr != nil {
		c.closeSocketConnection()
		return false, nil, sentErr
	}

	msg, rcvErr := c.waitResponse(signals)
	c.closeSocketConnection()

	return rcvErr == nil && bytesWritten == len(chunkData) && len(msg) > 0, msg, rcvErr
}
