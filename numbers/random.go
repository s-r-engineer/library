package libraryNumbers

import (
	"crypto/rand"
	"encoding/binary"
	"math/big"
)

func SimpleRand() (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

func RandomUint64() (uint64, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(b[:]), nil
}
