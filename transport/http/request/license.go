package request

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/repository/crypto"
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

	plaintext, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		return nil, err
	}

	var license LicenseRenew
	return &license, json.Unmarshal(plaintext, &license)
}
