package license

import (
	"crypto"
	"encoding/hex"
	"time"

	"github.com/pkg/errors"
	lic "github.com/ubogdan/license"

	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/repository"
	"github.com/ubogdan/network-manager-api/repository/serial"
	"github.com/ubogdan/network-manager-api/service"
)

type license struct {
	License         repository.License
	LicenseSigner   crypto.Signer
	SerialNumberKey []byte
}

var _ service.License = New(nil, nil, nil)

func New(lic repository.License, privateKey []byte, signer crypto.Signer) *license {
	return &license{
		License:         lic,
		LicenseSigner:   signer,
		SerialNumberKey: privateKey,
	}
}

func (s *license) FindAll() ([]model.License, error) {
	return s.License.FindAll()
}

func (s *license) Find(id uint64) (*model.License, error) {
	return s.License.Find(id)
}

func (s *license) Create(license *model.License) error {

	if license.Created == 0 {
		license.Created = time.Now().Unix()
	}

	if license.Expire == 0 {
		license.Expire = time.Unix(license.Created, 0).Add(12 * model.DefaultValidity).Unix()
	}

	validFromTime, err := nextValidPeriod(time.Unix(license.Created, 0), time.Unix(license.Expire, 0), time.Now(), model.DefaultValidity)
	if err != nil {
		return err
	}

	validUntilTime := validFromTime.Add(model.DefaultValidity) // 1 month

	licenseSerial, err := serial.Generate(hex.EncodeToString(s.SerialNumberKey), license.HardwareID, validUntilTime.Unix())
	if err != nil {
		return err
	}

	license.Serial = licenseSerial

	return s.License.Create(license)
}

func (s *license) Update(license *model.License) error {
	return s.License.Update(license)
}

func (s *license) Delete(id uint64) error {
	return s.License.Delete(id)
}

func (s *license) Renew(l *model.License) ([]byte, error) {

	license, err := s.License.FindByHardwareID(l.HardwareID)
	if err != nil {
		return nil, model.LicenseNotFound
	}

	if l.Serial != license.Serial {
		return nil, model.LicenseNotFound
	}

	//
	validFromTime, err := nextValidPeriod(time.Unix(license.Created, 0), time.Unix(license.Expire, 0), time.Now(), model.DefaultValidity)
	if err != nil {
		return nil, err
	}
	license.LastIssued = validFromTime.Unix()

	validUntilTime := validFromTime.Add(model.DefaultValidity) // 1 month

	licenseSerial, err := serial.Generate(hex.EncodeToString(s.SerialNumberKey), license.HardwareID, validUntilTime.Unix())
	if err != nil {
		return nil, err
	}

	// Create basic license Data
	licenseData := lic.License{
		ProductName:  model.ProductName,
		SerialNumber: licenseSerial,
		MinVersion:   10000,
		MaxVersion:   200000000,
		ValidFrom:    validFromTime,
		ValidUntil:   validUntilTime,
	}

	// Add license Features
	for _, feature := range license.Features {
		licenseData.Features = append(licenseData.Features, lic.Feature{
			Oid:         feature.Name.Oid(),
			Description: feature.Name.Description(),
			Expire:      feature.Expire,
			Limit:       feature.Limit,
		})
	}

	// Generate license
	lbytes, err := lic.CreateLicense(&licenseData, s.LicenseSigner)
	if err != nil {
		return nil, err
	}

	// Update database details
	err = s.License.Update(license)
	if err != nil {
		return nil, err
	}

	return lbytes, nil
}

func nextValidPeriod(created, expire, now time.Time, validity time.Duration) (time.Time, error) {
	if created.Add(validity).After(expire) {
		return created, errors.Wrapf(model.LicenseExpired, "no more days to add from %d to %d for cureent validity period of %d days", created.Unix(), expire.Unix(), validity/24/time.Hour)
	}

	if created.Add(validity).Before(now) {
		return nextValidPeriod(created.Add(validity), expire, now, validity)
	}

	return created, nil
}
