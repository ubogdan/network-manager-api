package repository

import (
	"github.com/ubogdan/network-manager-api/model"
)

// License repository definition.
type License interface {
	FindAll() ([]model.License, error)
	Find(hardwareID string) (*model.License, error)
	Create(license *model.License) error
	Update(license *model.License) error
	Delete(hardwareID string) error
}
