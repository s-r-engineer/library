package libraryEncryption

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"math/big"

	libraryErrors "github.com/s-r-engineer/library/errors"
	libraryNetwork "github.com/s-r-engineer/library/network"
)

func GetDHSecretFromConnection(conn libraryNetwork.GenericConnection, p *big.Int, g *big.Int) ([]byte, error) {
	wrappedError := libraryErrors.PartWrapError("GetDHSecretFromConnection")

	priv, err := rand.Int(rand.Reader, p)
	if err != nil {
		return nil, wrappedError(err)
	}

	pub := new(big.Int).Exp(g, priv, p)

	if _, err := conn.Write(pub.Bytes()); err != nil {
		return nil, wrappedError(err)
	}

	// Prepare a buffer big enough for the public key
	otherSidePub := make([]byte, p.BitLen()/8+1)

	// Read fully the peer's public key
	if _, err := io.ReadFull(conn, otherSidePub); err != nil {
		return nil, wrappedError(err)
	}

	otherSide := new(big.Int).SetBytes(otherSidePub)

	if otherSide.Cmp(big.NewInt(1)) <= 0 || otherSide.Cmp(new(big.Int).Sub(p, big.NewInt(1))) >= 0 {
		return nil, wrappedError(fmt.Errorf("invalid public key received"))
	}

	sharedSecret := new(big.Int).Exp(otherSide, priv, p)
	symmetricKey := sha256.Sum256(sharedSecret.Bytes())

	return symmetricKey[:], nil
}

func GetECDHSecretFromConnection(conn libraryNetwork.GenericConnection) ([]byte, error) {
	wrappedError := libraryErrors.PartWrapError("GetECDHSecretFromConnection")

	curve := ecdh.P521()

	privKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, wrappedError(err)
	}

	pubKeyBytes := privKey.PublicKey().Bytes()
	pubKeySize := len(pubKeyBytes)

	if _, err := conn.Write(pubKeyBytes); err != nil {
		return nil, wrappedError(err)
	}

	receivedBytes := make([]byte, pubKeySize)

	if _, err := io.ReadFull(conn, receivedBytes); err != nil {
		return nil, wrappedError(err)
	}

	otherPubKey, err := curve.NewPublicKey(receivedBytes)
	if err != nil {
		return nil, wrappedError(err)
	}

	sharedSecret, err := privKey.ECDH(otherPubKey)
	if err != nil {
		return nil, wrappedError(err)
	}

	symmetricKey := sha256.Sum256(sharedSecret)
	return symmetricKey[:], nil
}
