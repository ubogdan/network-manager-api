package response

import (
	"encoding/base64"
	"net/http"

	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/repository/crypto"
)

type License struct {
	Created    int64  `json:"created"`
	Expire     int64  `json:"expire"`
	Serial     string `json:"serial"`
	LastIssued int64  `json:"last_issued"`
}

func FromLicese(lic *model.License) License {
	return License{
		Created:    lic.Created,
		Expire:     lic.Expire,
		Serial:     lic.Serial,
		LastIssued: lic.LastIssued,
	}
}

func LicenseToEncryptedPayload(w http.ResponseWriter, payload, key []byte) error {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(200)

	writer := base64.NewEncoder(base64.StdEncoding, w)

	bytes, err := crypto.Encrypt(key, payload)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = writer.Write(bytes)

	return err
}
