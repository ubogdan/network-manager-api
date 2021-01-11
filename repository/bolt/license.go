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

func License(boltholdStore *bolthold.Store) *license {
	return &license{
		store: boltholdStore,
	}
}

func (s *license) FindAll() ([]model.License, error) {
	var licenses []model.License
	return licenses, s.store.Find(&licenses, bolthold.Where(bolthold.Key).Ne(uint64(0)))
}

func (s *license) Find(id uint64) (*model.License, error) {
	var license model.License
	return &license, s.store.FindOne(&license, bolthold.Where(bolthold.Key).Eq(id))
}

func (s *license) FindByHardwareID(hardwareID string) (*model.License, error) {
	var license model.License
	return &license, s.store.FindOne(&license, bolthold.Where("HardwareID").Eq(hardwareID))
}

func (s *license) Create(license *model.License) error {
	return s.store.Insert(bolthold.NextSequence(), license)
}

func (s *license) Update(license *model.License) error {
	return s.store.Update(license.ID, license)
}

func (s *license) Delete(id uint64) error {
	return s.store.Delete(id, model.License{})
}
