package request

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"github.com/ubogdan/network-manager-api/model"
)

type License struct {
	ID         uint64   `json:"id"`
	HardwareID string   `json:"hardware_id"`
	Customer   Customer `json:"customer"`
}

type LicenseRenew struct {
	Serial     string `json:"serial"`
	HardwareID string `json:"hardware_id"`
}

func (l *License) ToModel() model.License {
	return model.License{
		ID:         l.ID,
		HardwareID: l.HardwareID,
		Customer: model.Customer{
			Name:         l.Customer.Name,
			Country:      l.Customer.Country,
			City:         l.Customer.City,
			Organization: l.Customer.Organization,
		},
	}
}

func (l *LicenseRenew) ToModel() model.License {
	return model.License{
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

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesGCM.Open(nil, ciphertext[:12], ciphertext[12:], nil)
	if err != nil {
		return nil, err
	}

	var license LicenseRenew
	return &license, json.Unmarshal(plaintext, &license)
}
