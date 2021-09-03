package response

import (
	"encoding/base64"
	"net/http"

	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/repository/crypto"
)

// License response DTO.
type License struct {
	Created    int64     `json:"created"`
	Expire     int64     `json:"expire"`
	Serial     string    `json:"serial"`
	LastIssued int64     `json:"last_issued"`
	Customer   *Customer `json:"customer,omitempty"`
	Features   []Feature `json:"features,omitempty"`
}

type Feature struct {
	Name   string `json:"name"`
	Expire int64  `json:"expire"`
	Limit  int64  `json:"limit"`
}

// FromLicense convert license model to license DTO.
func FromLicense(lic *model.License) License {
	features := make([]Feature, 0, len(lic.Features))
	for _, feature := range lic.Features {
		features = append(features, Feature{
			Name:   string(feature.Name),
			Expire: feature.Expire,
			Limit:  feature.Limit,
		})
	}

	return License{
		Created:    lic.Created,
		Expire:     lic.Expire,
		Serial:     lic.Serial,
		LastIssued: lic.LastIssued,
		Customer: &Customer{
			Name:         lic.Customer.Name,
			Country:      lic.Customer.Country,
			City:         lic.Customer.City,
			Organization: lic.Customer.Organization,
		},
		Features: features,
	}
}

func FromLicenses(licenses []model.License) []License {
	response := make([]License, 0, len(licenses))
	for _, lic := range licenses {
		response = append(response, License{
			Created:    lic.Created,
			Expire:     lic.Expire,
			Serial:     lic.Serial,
			LastIssued: lic.LastIssued,
		})
	}
	return response
}

// LicenseToEncryptedPayload convert a license response to an encrypted payload.
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
