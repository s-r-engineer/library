package libraryEncryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"io"

	libraryErrors "github.com/s-r-engineer/library/errors"
	"golang.org/x/crypto/pbkdf2"
)

const (
	keyLength   = 32
	nonceLength = 12
	iterations  = 100000
)

func NewED(passphrase, salt string) (*ED, error) {
	nonce := make([]byte, nonceLength)
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}
	return &ED{passphrase: passphrase, salt: []byte(salt), nonce: nonce}, nil
}

type ED struct {
	salt []byte
	passphrase  string
	lock, unlock func()
	nonce int
}

func (e *ED) getNonce() []byte {

}

func (e *ED) deriveKey() []byte {
	return pbkdf2.Key([]byte(e.passphrase), e.salt, iterations, keyLength, sha512.New)
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
	if len(encryptedBytes) < nonceLength {
		return nil, libraryErrors.WrapError("ciphertext too short", nil)
	}
	nonce, ciphertext := encryptedBytes[:nonceLength], encryptedBytes[nonceLength:]

	return aesGCM.Open(nil, nonce, ciphertext, nil)
}
