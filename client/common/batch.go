package common

import (
	"archive/zip"
	"io"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// "slices"

type Batch struct {
	filePath      string
	agency_number int
}

func NewBatch(clientConfig ClientConfig) *Batch {

	agency_number, err := strconv.Atoi(clientConfig.ID)

	if err != nil {
		return nil
	}

	return &Batch{
		filePath:      clientConfig.ZipBatchPath,
		agency_number: agency_number,
	}
}

func (batch *Batch) recursivelyReadFile(zipReader io.ReadCloser, agencyFileName string, chunk_number int) ([]byte, error) {

	defer zipReader.Close()
	buffer := make([]byte, 1024)
	n, readErr := zipReader.Read(buffer)

	if readErr != nil && readErr != io.EOF {
		log.Errorf("action: reading_file | result: fail | file: %v | chunk_number: %v | bytes_len: %v", agencyFileName, chunk_number, n)
		return nil, readErr

	} else if n < len(buffer) || readErr == io.EOF {
		return buffer[:n], nil

	} else {
		data, err := batch.recursivelyReadFile(zipReader, agencyFileName, chunk_number+1)

		if err != nil && err != io.EOF {
			log.Infof("action: building_content | result: fail | file: %v | bytes_len: %v", agencyFileName, err)
			return nil, err
		}

		content := append(buffer, data...)

		return content, nil
	}
}

func (batch *Batch) ReadAgencyFileContents(zipReader *zip.ReadCloser, ch_msg chan []byte, ch_err chan error) {

	agencyFileName := "agency-" + strconv.Itoa(batch.agency_number) + ".csv"
	for _, file := range zipReader.File {
		if file.FileHeader.Name == agencyFileName {
			log.Infof("action: zip file opened | filename: %v", file.FileHeader.Name)

			ioReadCloser, agencyFileError := file.Open()

			if agencyFileError != nil {
				log.Infof("action: csv_file_open | result: fail | client_id: %v | file_name: %v | msg: %v",
					batch.agency_number,
					agencyFileName,
					agencyFileError,
				)
				ch_err <- agencyFileError
				break
			}

			defer ioReadCloser.Close()

			fileBytes, zipContentError := batch.recursivelyReadFile(ioReadCloser, agencyFileName, 1)

			if zipContentError == nil || zipContentError == io.EOF {
				log.Infof("action: csv_file_processed | result: success | client_id: %v | file_name: %v | file_size: %v",
					batch.agency_number,
					agencyFileName,
					len(fileBytes),
				)

				ch_msg <- fileBytes

			} else {
				log.Errorf("action: csv_file_processed | result: fail | client_id: %v | file_name: %v | err: %v",
					batch.agency_number,
					agencyFileName,
					zipContentError,
				)
				ch_err <- zipContentError
			}
			break
		}
	}

}

func (batch *Batch) ProcessBatchFile() ([]byte, error) {
	ch_sig := make(chan os.Signal, 1)
	ch_msg := make(chan []byte)
	ch_err := make(chan error)
	signal.Notify(ch_sig, syscall.SIGTERM)

	zipFile, err := zip.OpenReader(batch.filePath)

	if err != nil {
		log.Infof("action: unzipping_file | result: fail | client_id: %v | file_name: %v | err: %v",
			batch.agency_number,
			batch.filePath,
			err,
		)
	}

	defer zipFile.Close()

	go batch.ReadAgencyFileContents(zipFile, ch_msg, ch_err)

	select {
	case msg := <-ch_msg:
		log.Infof("action: file_unzipped | result: succeed | client_id: %v | file_name: %v | msg_size: %v",
			batch.agency_number,
			batch.filePath,
			len(msg),
		)
		return msg, nil

	case err := <-ch_err:
		log.Infof("action: file_unzipped | result: fail | client_id: %v | file_name: %v | err: %v",
			batch.agency_number,
			batch.filePath,
			err,
		)
		return nil, err

	case signal := <-ch_sig:
		log.Infof("action: sigterm_signal | result: success | client_id: %v | msg: %v signal received",
			batch.agency_number,
			signal,
		)
		return nil, nil
	}

}
