package request

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	models "github.com/ubogdan/network-manager-api/model"
)

type License struct {
	ID         int64  `json:"id"`
	HardwareID string `json:"hardware_id"`
	Customer   string `json:"customer"`
}

type LicenseRenew struct {
	Serial     string `json:"serial"`
	HardwareID string `json:"hardware_id"`
}

func (l *License) ToModel() models.License {
	return models.License{
		ID:         l.ID,
		HardwareID: l.HardwareID,
		Customer:   l.Customer,
	}
}

func (l *LicenseRenew) ToModel() models.License {
	return models.License{
		HardwareID: l.HardwareID,
		Serial:     l.Serial,
	}
}

func LicenseFromEncryptedPayload(reader io.ReadCloser, key []byte) (*LicenseRenew, error) {
	ciphertext, err := ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, reader))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	if len(ciphertext) < aes.BlockSize+12 {
		return nil, errors.New("ciphertext too short")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, ciphertext[:12], ciphertext[12:], nil)
	if err != nil {
		return nil, err
	}

	var license LicenseRenew
	return &license, json.Unmarshal(plaintext, &license)
}
