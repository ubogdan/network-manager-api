package repository

import (
	"github.com/ubogdan/network-manager-api/model"
)

type License interface {
	FindAll() ([]model.License, error)
	Find(id uint64) (*model.License, error)
	Create(license *model.License) error
	Update(license *model.License) error
	Delete(id uint64) error
}
