package libraryNetwork

import (
	"io"

	libraryErrors "github.com/s-r-engineer/library/errors"
)

func ReadConnection(conn GenericConnection, length int) (data []byte, err error) {
	bufferLength := 1024
	if length > 0 {
		bufferLength = length
	}
	buffer := make([]byte, bufferLength)
	for {
		n, err := conn.Read(buffer)
		if n > 0 {
			data = append(data, buffer[:n]...)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return data, err
		}
		if n < bufferLength || len(buffer) == length {
			break
		}
	}
	return
}

func WriteConnection(conn GenericConnection, data []byte) error {
	total := 0
	for total < len(data) {
		n, err := conn.Write(data[total:])
		if err != nil {
			return err
		}
		total += n
	}
	return nil
}

func ConnectPipes(conn1, conn2 GenericConnection, errChan chan error) {
	copyData := func(dst, src GenericConnection, errChan chan error) {
		_, err := io.Copy(dst, src)
		if errChan != nil {
			errChan <- err
		} else {
			libraryErrors.Errorer(err)
		}
	}
	go copyData(conn1, conn2, errChan)
	go copyData(conn2, conn1, errChan)
}
