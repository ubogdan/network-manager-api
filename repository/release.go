package repository

import (
	models "github.com/ubogdan/network-manager-api/model"
)

type Release interface {
	Find(channel string) (*models.Release, error)
	Create(release *models.Release) error
	Update(release *models.Release) error
}
