package bolt

import (
	models "github.com/ubogdan/network-manager-api/model"
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

func (s *license) FindAll() ([]models.License, error) {
	return nil, nil
}

func (s *license) Find(id int64) (*models.License, error) {
	return nil, nil
}

func (s *license) Create(license *models.License) error {
	return nil
}

func (s *license) Update(license *models.License) error {
	return nil
}

func (s *license) Delete(id int64) error {
	return nil
}
