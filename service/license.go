package service

import (
	"github.com/ubogdan/network-manager-api/model"
)

// License service definition.
type License interface {
	FindAll() ([]model.License, error)
	Find(id uint64) (*model.License, error)
	FindByHardwareID(hardwareID string) (*model.License, error)
	Create(license *model.License) error
	Update(license *model.License) error
	Delete(id uint64) error
	Renew(license *model.License) ([]byte, error)
}
