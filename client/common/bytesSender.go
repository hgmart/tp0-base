package common

import (
	"os"
	"os/signal"
	"syscall"
)

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

func (c *Client) SendBytes(chunkData []byte) (bool, error) {
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

	bytesWritten, err := c.SendData(chunkData, signals, 1)
	log.Infof("action: mensaje_enviado | result: %v | client: %v | chunk_size: %v | bytes_sent: %v", err == nil, c.config.ID, len(chunkData), bytesWritten)

	c.closeSocketConnection()

	return err == nil && bytesWritten == len(chunkData), err
}
