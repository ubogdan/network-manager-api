package repository

import (
	models "github.com/ubogdan/network-manager-api/model"
)

type License interface {
	FindAll() ([]models.License, error)
	Find(id int64) (*models.License, error)
	Create(license *models.License) error
	Update(license *models.License) error
	Delete(id int64) error
}
