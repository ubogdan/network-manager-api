package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// Encrypt a payload with the given key using aesGCM cipher.
func Encrypt(key, payload []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	seal, err := aesGCM.Seal(nil, nonce, payload, nil), nil
	if err != nil {
		return nil, err
	}

	return append(nonce, seal...), nil
}

// EncryptWithStringKey is helper function what can be used with hex encoded key.
func EncryptWithStringKey(privateKey string, payload []byte) ([]byte, error) {
	key, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}

	return Encrypt(key, payload)
}

// Decrypt an aesGCM payload using a given key.
func Decrypt(key, payload []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesGCM.Open(nil, payload[:aesGCM.NonceSize()], payload[aesGCM.NonceSize():], nil)
}

// DecryptWithStringKey is helper function what can be used with hex encoded key.
func DecryptWithStringKey(privateKey string, payload []byte) ([]byte, error) {
	key, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	return Decrypt(key, payload)
}
