package libraryEncryption

import (
	"testing"

	"github.com/stretchr/testify/require"

	libraryStrings "github.com/s-r-engineer/library/strings"
)

func TestEncryptor(t *testing.T) {
	salt := libraryStrings.RandString(16)
	data := libraryStrings.RandString(666)
	passphrase := libraryStrings.RandString(32)
	encrypted, err := EncryptAES(passphrase, salt, []byte(data))
	require.NoError(t, err)
	decrypted, err := DecryptAES(passphrase, salt, encrypted)
	require.NoError(t, err)
	require.Equal(t, data, string(decrypted))
}

func BenchmarkEncryptor(b *testing.B) {
	salt := libraryStrings.RandString(16)
	data := libraryStrings.RandString(666)
	passphrase := libraryStrings.RandString(32)
	for b.Loop() {
		encrypted, _ := EncryptAES(passphrase, salt, []byte(data))
		_, _ = DecryptAES(passphrase, salt, encrypted)

	}
}
