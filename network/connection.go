package libraryNetwork

import (
	"io"
	"net"

	libraryErrors "github.com/s-r-engineer/library/errors"
)

func ReadConnection(conn net.Conn) (data []byte, err error) {
	n := 1024
	buffer := make([]byte, n)
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
	}
	return
}

func WriteConnection(conn net.Conn, data []byte) error {
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

func ConnectPipes(conn1, conn2 net.Conn) {
	copyData := func(dst, src net.Conn) {
		defer dst.Close()
		defer src.Close()
		_, err := io.Copy(dst, src)
		libraryErrors.Errorer(err)
	}
	go copyData(conn1, conn2)
	go copyData(conn2, conn1)
}

