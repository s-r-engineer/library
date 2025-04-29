package hell

import (
	"encoding/binary"
	"io"
	"math/big"

	libraryEncryption "github.com/s-r-engineer/library/encryption"
	libraryErrors "github.com/s-r-engineer/library/errors"
	libraryNetwork "github.com/s-r-engineer/library/network"
)

func MakeAHellCircle(connection libraryNetwork.GenericConnection, salt string, p, g *big.Int) libraryNetwork.GenericConnection {
	secret, err := libraryEncryption.GetDHSecretFromConnection(connection, p, g)
	libraryErrors.Errorer(err)
	return &HellCircle{connection: connection, salt: salt, password: secret}
}

type HellCircle struct {
	connection libraryNetwork.GenericConnection
	salt       string
	password   string
}

func (m HellCircle) Read(b []byte) (int, error) {
	var len1 uint32
	err := binary.Read(m.connection, binary.BigEndian, &len1)
	if err != nil {
		return 0, err
	}
	var data = make([]byte, len1)
	_, err = io.ReadFull(m.connection, data)
	if err != nil {
		return 0, err
	}
	decryptedData, err := libraryEncryption.DecryptAES(m.password, m.salt, data)
	if err != nil {
		return 0, err
	}
	if len(b) < len(decryptedData) {
		return 0, io.ErrShortBuffer
	}
	n := copy(b, decryptedData)
	return n, nil
}

func (m HellCircle) Write(b []byte) (int, error) {
	encryptedData, err := libraryEncryption.EncryptAES(m.password, m.salt, b)
	if err != nil {
		return 0, err
	}
	len1 := uint32(len(encryptedData))
	err = binary.Write(m.connection, binary.BigEndian, len1)
	if err != nil {
		return 0, err
	}
	// m.connection.Write(encryptedData)
	_, err = m.connection.Write(encryptedData)
	if err != nil {
		return 0, err
	}
	return len(b), err
}

func (m HellCircle) Close() error {
	return m.connection.Close()
}
