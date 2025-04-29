package libraryEncryption

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	libraryNetwork "github.com/s-r-engineer/library/network"
)

func GetDHSecretFromConnection(conn libraryNetwork.GenericConnection, p *big.Int, g *big.Int) (string, error) {
	priv, err := rand.Int(rand.Reader, p)
	if err != nil {
		return "", err
	}
	pub := new(big.Int).Exp(g, priv, p)
	otherSidePub := make([]byte, p.BitLen()/8+1)
	_, err = conn.Write(pub.Bytes())
	if err != nil {
		return "", err
	}
	n, err := conn.Read(otherSidePub)
	if err != nil {
		return "", err
	}
	otherSide := new(big.Int).SetBytes(otherSidePub[:n])
	sharedSecret := new(big.Int).Exp(otherSide, priv, p)
	symmetricKey := sha256.Sum256(sharedSecret.Bytes())
	return fmt.Sprintf("%x", symmetricKey), nil
}

func GetECDHSecretFromConnection(conn libraryNetwork.GenericConnection) (string, error) {
	curve := ecdh.P521()
	privKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return "", err
	}
	pubKeyBytes := privKey.PublicKey().Bytes()
	pubKeySize := len(pubKeyBytes)

	if _, err = conn.Write(pubKeyBytes); err != nil {
		return "", err
	}

	otherPubKeyBytes := make([]byte, pubKeySize)
	n, err := conn.Read(otherPubKeyBytes)
	if err != nil {
		return "", err
	}
	if n != pubKeySize {
		return "", fmt.Errorf("received incorrect public key length: expected %d, got %d", pubKeySize, n)
	}

	otherPubKey, err := curve.NewPublicKey(otherPubKeyBytes)
	if err != nil {
		return "", fmt.Errorf("invalid public key received: %v", err)
	}

	sharedSecret, err := privKey.ECDH(otherPubKey)
	if err != nil {
		return "", err
	}

	symmetricKey := sha256.Sum256(sharedSecret)
	return fmt.Sprintf("%x", symmetricKey), nil
}
