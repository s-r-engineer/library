package libraryEncryption

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	libraryErrors "github.com/s-r-engineer/library/errors"
	libraryNetwork "github.com/s-r-engineer/library/network"
)

func GetDHSecretFromConnection(conn libraryNetwork.GenericConnection, p *big.Int, g *big.Int) (string, error) {
	wrappedError := libraryErrors.PartWrapError("GetDHSecretFromConnection")
	priv, err := rand.Int(rand.Reader, p)
	if err != nil {
		return "", wrappedError(err)
	}

	pub := new(big.Int).Exp(g, priv, p)

	if _, err := conn.Write(pub.Bytes()); err != nil {
		return "", wrappedError(err)
	}

	otherSidePub := make([]byte, p.BitLen()/8+1)
	n, err := conn.Read(otherSidePub)
	if err != nil {
		return "", wrappedError(err)
	}

	otherSide := new(big.Int).SetBytes(otherSidePub[:n])

	if otherSide.Cmp(big.NewInt(1)) <= 0 || otherSide.Cmp(new(big.Int).Sub(p, big.NewInt(1))) >= 0 {
		return "", wrappedError(err)
	}

	sharedSecret := new(big.Int).Exp(otherSide, priv, p)
	symmetricKey := sha256.Sum256(sharedSecret.Bytes())

	return fmt.Sprintf("%x", symmetricKey), nil
}

func GetECDHSecretFromConnection(conn libraryNetwork.GenericConnection) (string, error) {
	wrappedError, wrappedString := libraryErrors.PartWrapErrorOrString("GetDHSecretFromConnection")
	curve := ecdh.P521()

	privKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return "", wrappedError(err)
	}

	pubKeyBytes := privKey.PublicKey().Bytes()
	pubKeySize := len(pubKeyBytes)

	if _, err := conn.Write(pubKeyBytes); err != nil {
		return "", wrappedError(err)
	}

	receivedBytes := make([]byte, pubKeySize)
	n, err := conn.Read(receivedBytes)
	if err != nil {
		return "", wrappedError(err)
	}
	if n != pubKeySize {
		return "", wrappedString("received incorrect public key length")
	}

	otherPubKey, err := curve.NewPublicKey(receivedBytes[:n])
	if err != nil {
		return "", wrappedError(err)
	}

	sharedSecret, err := privKey.ECDH(otherPubKey)
	if err != nil {
		return "", wrappedError(err)
	}

	symmetricKey := sha256.Sum256(sharedSecret)
	return fmt.Sprintf("%x", symmetricKey), nil
}
