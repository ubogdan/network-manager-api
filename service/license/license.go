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
)

type license struct {
	License         repository.License
	LicenseSigner   crypto.Signer
	SerialNumberKey []byte
}

// New returns license service implementation.
func New(lic repository.License, privateKey []byte, signer crypto.Signer) *license {
	return &license{
		License:         lic,
		LicenseSigner:   signer,
		SerialNumberKey: privateKey,
	}
}

// FindAll returns license list.
func (s *license) FindAll() ([]model.License, error) {
	return s.License.FindAll()
}

// Find return license by id.
func (s *license) Find(id string) (*model.License, error) {
	return s.License.Find(id)
}

// Create new license.
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

// Update existing license.
func (s *license) Update(license *model.License) error {
	return s.License.Update(license)
}

// Delete license by id.
func (s *license) Delete(id string) error {
	return s.License.Delete(id)
}

// Renew license.
func (s *license) Renew(license *model.License) ([]byte, error) {
	validFromTime, err := nextValidPeriod(time.Unix(license.Created, 0), time.Unix(license.Expire, 0), time.Now().Add(model.DefaultGracePeriod), model.DefaultValidity)
	if err != nil {
		return nil, err
	}

	validUntilTime := validFromTime.Add(model.DefaultValidity) // 1 month

	if validFromTime.After(time.Now()) {
		validFromTime = time.Now()
	}

	license.LastIssued = validFromTime.Unix()

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

	// Update last serial
	license.Serial = licenseSerial

	// Update database details
	err = s.License.Update(license)
	if err != nil {
		return nil, err
	}

	return lbytes, nil
}

func nextValidPeriod(created, expire, now time.Time, validity time.Duration) (time.Time, error) {
	if created.Add(validity).After(expire) {
		return created, errors.Wrapf(model.ErrLicenseExpired, "no more days to add from %d to %d for cureent validity period of %d days", created.Unix(), expire.Unix(), validity/24/time.Hour)
	}

	if created.Add(validity).Before(now) {
		return nextValidPeriod(created.Add(validity), expire, now, validity)
	}

	return created, nil
}
