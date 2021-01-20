package response

import (
	"encoding/base64"
	"net/http"

	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/repository/crypto"
)

type License struct {
}

func FromLicese(lic *model.License) License {
	return License{}
}

func LicenseToEncryptedPayload(w http.ResponseWriter, payload, key []byte) error {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(200)

	bytes, err := crypto.Encrypt(key, payload)
	if err != nil {
		return err
	}

	_, err = base64.NewEncoder(base64.StdEncoding, w).Write(bytes)
	return err
}
