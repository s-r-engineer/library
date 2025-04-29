package libraryEncryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/binary"

	libraryErrors "github.com/s-r-engineer/library/errors"
	libraryNumbers "github.com/s-r-engineer/library/numbers"
	"golang.org/x/crypto/pbkdf2"
)

const (
	keyLength   = 32
	nonceLength = 12
	iterations  = 100000
)

func ValueOrDefault[T int | int64 | uint | uint64](value, defaults T, f func(T, T) bool) T {
	if value <= 0 || (f != nil && !f(value, defaults)) {
		return defaults
	}
	return value
}

func NewED(passphrase, salt string, keyLengthUser, iterationsUser, nonceLengthUser int) (*ED, error) {
	randomNonce, err := libraryNumbers.RandomUint64()
	if err != nil {
		return nil, err
	}
	return &ED{
		keyLength:   ValueOrDefault(keyLengthUser, keyLength, nil),
		iterations:  ValueOrDefault(iterationsUser, iterations, nil),
		nonceLength: ValueOrDefault(nonceLengthUser, nonceLength, nil),
		passphrase:  passphrase,
		salt:        []byte(salt),
		nonce:       randomNonce,
	}, nil
}

type ED struct {
	salt                               []byte
	passphrase                         string
	nonce                              uint64
	keyLength, nonceLength, iterations int
}

func (e *ED) getNonce() []byte {
	bytes := make([]byte, e.nonceLength)
	binary.BigEndian.PutUint64(bytes[4:], e.nonce)
	e.nonce++
	return bytes
}

func (e *ED) deriveKey() []byte {
	return pbkdf2.Key([]byte(e.passphrase), e.salt, e.iterations, e.keyLength, sha512.New)
}

func (e *ED) EncryptAES(plaintextBytes []byte) ([]byte, error) {
	key := e.deriveKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := e.getNonce()
	return aesGCM.Seal(nonce, nonce, plaintextBytes, nil), nil
}

func (e *ED) DecryptAES(encryptedBytes []byte) ([]byte, error) {
	key := e.deriveKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(encryptedBytes) < e.nonceLength {
		return nil, libraryErrors.WrapError("ciphertext too short", nil)
	}
	nonce, ciphertext := encryptedBytes[:e.nonceLength], encryptedBytes[e.nonceLength:]

	return aesGCM.Open(nil, nonce, ciphertext, nil)
}
