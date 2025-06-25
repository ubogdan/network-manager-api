package serial

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ubogdan/network-manager-api/repository/crypto"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	// Generate random private key
	privateKeyBytes := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, privateKeyBytes)
	assert.NoError(t, err)

	encrypt, err := crypto.Encrypt(privateKeyBytes, []byte(`
package serial

func Generate(key string, hwID string, valid int64) (string,error) {
	return "8c658cf2-FAKE-FAKE-FAKE-SERIAL130002", nil
}

`))
	assert.NoError(t, err)

	pemBlock.Bytes = encrypt

	result, err := Generate(hex.EncodeToString(privateKeyBytes), "1234567890", 0)
	assert.NoError(t, err)
	assert.Equal(t, result, "8c658cf2-FAKE-FAKE-FAKE-SERIAL130002")
}

func TestValidUntil(t *testing.T) {
	t.Parallel()

	// Generate random private key
	privateKeyBytes := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, privateKeyBytes)
	assert.NoError(t, err)

	encrypt, err := crypto.Encrypt(privateKeyBytes, []byte(`
package serial
import (
	"errors"
)

func ValidUntil(serial string) (int64,error) {
	if serial == "8c658cf2-FAKE-FAKE-FAKE-SERIAL130002" {
		return 1751957281, nil
	}
	return 0, errors.New("invalid serial")
}

`))
	assert.NoError(t, err)

	pemValidUntil.Bytes = encrypt

	result, err := ValidUntil(hex.EncodeToString(privateKeyBytes), "8c658cf2-FAKE-FAKE-FAKE-SERIAL130002")
	assert.NoError(t, err)
	assert.Equal(t, result, int64(1751957281))
}
