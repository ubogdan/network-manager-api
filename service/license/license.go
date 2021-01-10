package license

import (
	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/repository"
	"github.com/ubogdan/network-manager-api/service"
)

type license struct {
	License repository.License
}

var _ service.License = New(nil)

func New(lic repository.License) *license {
	return &license{
		License: lic,
	}
}

func (s *license) FindAll() ([]model.License, error) {
	return nil, nil
}

func (s *license) Find(id int64) (*model.License, error) {
	return nil, nil
}

func (s *license) Create(license *model.License) error {
	return nil
}

func (s *license) Update(license *model.License) error {
	return nil
}

func (s *license) Delete(id int64) error {
	return nil
}

func (s *license) Renew(license *model.License) error {
	return nil
}
