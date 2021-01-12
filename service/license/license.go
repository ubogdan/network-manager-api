package license

import (
	"crypto"
	"time"

	lic "github.com/ubogdan/license"

	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/repository"
	"github.com/ubogdan/network-manager-api/repository/serial"
	"github.com/ubogdan/network-manager-api/service"
)

type license struct {
	Signer  crypto.Signer
	License repository.License
}

var _ service.License = New(nil, nil)

func New(lic repository.License, key crypto.Signer) *license {
	return &license{
		License: lic,
	}
}

func (s *license) FindAll() ([]model.License, error) {
	return s.License.FindAll()
}

func (s *license) Find(id uint64) (*model.License, error) {
	return s.License.Find(id)
}

func (s *license) Create(license *model.License) error {
	return s.Create(license)
}

func (s *license) Update(license *model.License) error {
	return s.Update(license)
}

func (s *license) Delete(id uint64) error {
	return s.License.Delete(id)
}

func (s *license) Renew(license *model.License) ([]byte, error) {

	_, err := s.License.FindByHardwareID(license.HardwareID)
	if err != nil {
		return nil, model.LicenseNotFound
	}

	licenseSerial, err := serial.Generate("", license.HardwareID, 1)
	if err != nil {
		return nil, err
	}

	// TODO : from license data
	validFromTime := time.Now()
	validUntilTime := validFromTime.Add(30 * 24 * time.Hour) // 1 month

	licenseData := lic.License{
		ProductName:  model.ProductName,
		SerialNumber: licenseSerial,
		MinVersion:   10000,
		MaxVersion:   200000000,
		ValidFrom:    validFromTime,
		ValidUntil:   validUntilTime,
	}

	return lic.CreateLicense(&licenseData, s.Signer)
}
