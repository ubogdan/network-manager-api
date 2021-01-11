package response

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/ubogdan/network-manager-api/model"
)

type License struct {
}

func FromLicese(lic *model.License) License {
	return License{}
}

func LicenseToEncryptedPayload(w http.ResponseWriter, payload, key []byte) error {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(200)

	writer := base64.NewEncoder(base64.StdEncoding, w)
	defer writer.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	_, err = writer.Write(aesGCM.Seal(nil, nonce, payload, nil))
	return err
}
