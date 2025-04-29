package libraryEncryption_test

import (
	"bytes"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	libraryEncryption "github.com/s-r-engineer/library/encryption"
)

type LinkedMockConnection struct {
	peer   *LinkedMockConnection
	buffer *bytes.Buffer
	mu     sync.Mutex
	closed bool
}

func NewLinkedMockConnections() (*LinkedMockConnection, *LinkedMockConnection) {
	c1 := &LinkedMockConnection{}
	c2 := &LinkedMockConnection{}
	c1.buffer = new(bytes.Buffer)
	c2.buffer = new(bytes.Buffer)
	c1.peer = c2
	c2.peer = c1
	return c1, c2
}

func (c *LinkedMockConnection) Read(p []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for {
		if c.closed {
			return 0, errors.New("connection closed")
		}
		if c.buffer.Len() > 0 {
			return c.buffer.Read(p)
		}
		c.mu.Unlock()
		time.Sleep(10 * time.Millisecond)
		c.mu.Lock()
	}
}

func (c *LinkedMockConnection) Write(p []byte) (int, error) {
	c.peer.mu.Lock()
	defer c.peer.mu.Unlock()
	if c.closed || c.peer.closed {
		return 0, errors.New("connection closed")
	}
	return c.peer.buffer.Write(p)
}

func (c *LinkedMockConnection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	return nil
}

func TestGetDHSecretFromConnection(t *testing.T) {
	p, _ := new(big.Int).SetString("23", 10)
	g := big.NewInt(5)

	clientConn, serverConn := NewLinkedMockConnections()

	var clientSecret, serverSecret string
	var clientErr, serverErr error

	done := make(chan struct{})

	go func() {
		clientSecret, clientErr = libraryEncryption.GetDHSecretFromConnection(clientConn, p, g)
		done <- struct{}{}
	}()

	go func() {
		serverSecret, serverErr = libraryEncryption.GetDHSecretFromConnection(serverConn, p, g)
		done <- struct{}{}
	}()

	<-done
	<-done

	if clientErr != nil || serverErr != nil {
		t.Fatalf("errors: clientErr=%v serverErr=%v", clientErr, serverErr)
	}

	if clientSecret != serverSecret {
		t.Fatalf("shared secrets do not match: client=%s server=%s", clientSecret, serverSecret)
	}
}



func TestGetECDHSecretFromConnection(t *testing.T) {
	clientConn, serverConn := NewLinkedMockConnections()

	var clientSecret, serverSecret string
	var clientErr, serverErr error

	done := make(chan struct{})

	go func() {
		clientSecret, clientErr = libraryEncryption.GetECDHSecretFromConnection(clientConn)
		done <- struct{}{}
	}()

	go func() {
		serverSecret, serverErr = libraryEncryption.GetECDHSecretFromConnection(serverConn)
		done <- struct{}{}
	}()

	<-done
	<-done

	if clientErr != nil || serverErr != nil {
		t.Fatalf("errors: clientErr=%v serverErr=%v", clientErr, serverErr)
	}

	if clientSecret != serverSecret {
		t.Fatalf("shared secrets do not match: client=%s server=%s", clientSecret, serverSecret)
	}
}
