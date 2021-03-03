package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	// Generate random private key
	privateKeyBytes := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, privateKeyBytes)

	samplePayload := []byte(`Super Secret Message`)
	encryptedBytes, err := Encrypt(privateKeyBytes, samplePayload)
	assert.NoError(t, err)

	expect, err := Decrypt(privateKeyBytes, encryptedBytes)
	assert.NoError(t, err)
	assert.Equal(t, expect, samplePayload)
}

func TestEncryptDecryptWithString(t *testing.T) {
	t.Parallel()

	// Generate random private key
	privateKeyBytes := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, privateKeyBytes)
	assert.NoError(t, err)

	privateKey := hex.EncodeToString(privateKeyBytes)

	samplePayload := []byte(`Super Secret Message`)
	encr, err := EncryptWithStringKey(privateKey, samplePayload)
	assert.NoError(t, err)

	expect, err := DecryptWithStringKey(privateKey, encr)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, expect, samplePayload)
}
