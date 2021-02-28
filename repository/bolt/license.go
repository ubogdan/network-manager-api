package bolt

import (
	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/pkg/bolthold"
	"github.com/ubogdan/network-manager-api/repository"
)

type license struct {
	store *bolthold.Store
}

var _ repository.License = License(nil)

// License return a license repository
func License(boltholdStore *bolthold.Store) *license {
	return &license{
		store: boltholdStore,
	}
}

// FindAll returns a list of licenses
func (s *license) FindAll() ([]model.License, error) {
	var licenses []model.License
	return licenses, s.store.Find(&licenses, bolthold.Where(bolthold.Key).Ne(uint64(0)))
}

// Find returns a license by id
func (s *license) Find(id uint64) (*model.License, error) {
	var license model.License
	return &license, s.store.FindOne(&license, bolthold.Where(bolthold.Key).Eq(id))
}

// FindByHardwareID returns a license by HardwareID
func (s *license) FindByHardwareID(hardwareID string) (*model.License, error) {
	var license model.License
	return &license, s.store.FindOne(&license, bolthold.Where("HardwareID").Eq(hardwareID))
}

// Create a new license record
func (s *license) Create(license *model.License) error {
	var lic model.License
	err := s.store.FindOne(&lic, bolthold.Where("HardwareID").Eq(license.HardwareID))
	if err != nil {
		return s.store.Insert(bolthold.NextSequence(), license)
	}
	return model.LicenseAlreadyExists
}

// Update a license record
func (s *license) Update(license *model.License) error {
	return s.store.Update(license.ID, license)
}

// Delete a license record
func (s *license) Delete(id uint64) error {
	return s.store.Delete(id, model.License{})
}
