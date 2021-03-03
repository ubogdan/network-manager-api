package request

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/repository/crypto"
)

// License request DTO.
type License struct {
	ID         uint64   `json:"id"`
	Created    int64    `json:"created"`
	Expire     int64    `json:"expire"`
	HardwareID string   `json:"hardware_id"`
	Customer   Customer `json:"customer"`
}

// Renew request DTO.
type Renew struct {
	Serial     string `json:"serial"`
	HardwareID string `json:"hardware_id"`
}

// ToModel converts license DTO to license model.
func (l *License) ToModel() model.License {
	return model.License{
		ID:         l.ID,
		Created:    l.Created,
		Expire:     l.Expire,
		HardwareID: l.HardwareID,
		Customer: model.Customer{
			Name:         l.Customer.Name,
			Country:      l.Customer.Country,
			City:         l.Customer.City,
			Organization: l.Customer.Organization,
		},
	}
}

// ToModel converts renew DTO to license model.
func (l *Renew) ToModel() model.License {
	return model.License{
		HardwareID: l.HardwareID,
		Serial:     l.Serial,
	}
}

// LicenseFromEncryptedPayload godoc.
func LicenseFromEncryptedPayload(reader io.ReadCloser, key []byte) (*Renew, error) {
	ciphertext, err := ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, reader))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	plaintext, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		return nil, err
	}

	var license Renew

	return &license, json.Unmarshal(plaintext, &license)
}
