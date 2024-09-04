package communication

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

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
