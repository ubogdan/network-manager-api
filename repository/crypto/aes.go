package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

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
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encrBytes, err := aesGCM.Seal(nil, nonce, payload, nil), nil
	if err != nil {
		return nil, err
	}
	return append(nonce, encrBytes...), nil
}

func EncryptWithStringKey(privateKey string, payload []byte) ([]byte, error) {
	key, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	return Encrypt(key, payload)
}

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

func DecryptWithStringKey(privateKey string, payload []byte) ([]byte, error) {
	key, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	return Decrypt(key, payload)
}
