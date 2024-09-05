package communication

import (
	"bytes"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

func BatchProcessing(chunkType string, agency string, data []byte) []byte {

	propertyDivider := []byte{0}
	chunksDivider := []byte{1}

	block0 := []byte(chunkType)
	block1 := []byte(agency)

	// Modifica el separador \r\n por el byte 1
	newdata := bytes.ReplaceAll(data, []byte("\r\n"), chunksDivider)

	// Modifica el separador ASCII ',' por el byte 0
	newdata = bytes.ReplaceAll(newdata, []byte{44}, propertyDivider)

	messageType := append(block0, propertyDivider...)
	messageAgency := append(block1, chunksDivider...)

	messageMetadata := append(messageType, messageAgency...)

	return append(messageMetadata, newdata...)
}

func Build(chunkType string, data []string) []byte {

	block0 := []byte(chunkType)
	separator := []byte{0}

	log.Infof("action: block_0    | result: success | size: %v", len(block0))
	log.Infof("action: separator  | result: success | size: %v", len(separator))

	result := append(block0, separator...)

	for index := 0; index < len(data); index++ {
		bytesX := []byte(data[index])
		log.Infof("action: block_%v    | result: success | size: %v", index+1, len(bytesX))
		log.Infof("action: separator  | result: success | size: %v", len(separator))

		result = append(result, append(bytesX, separator...)...)
	}

	log.Infof("action: chunk_size | result: success | size: %v", len(result))

	return result
}
