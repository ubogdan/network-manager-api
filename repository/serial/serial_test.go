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

	// Generate random private key
	privateKeyBytes := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, privateKeyBytes)
	assert.NoError(t, err)

	encr, err := crypto.Encrypt(privateKeyBytes, []byte(`
package serial

func Generate(key string, hwID string, valid int64) (string,error) {
	return "8c658cf2-FAKE-FAKE-FAKE-SERIAL130002", nil
}

`))
	assert.NoError(t, err)

	pemBlock.Bytes = encr

	result, err := Generate(hex.EncodeToString(privateKeyBytes), "1234567890", 0)
	assert.NoError(t, err)
	assert.Equal(t, result, "8c658cf2-FAKE-FAKE-FAKE-SERIAL130002")

}
