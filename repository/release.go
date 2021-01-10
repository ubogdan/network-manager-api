package repository

import (
	"github.com/ubogdan/network-manager-api/model"
)

type Release interface {
	Find(channel string) (*model.Release, error)
	Create(release *model.Release) error
	Update(release *model.Release) error
}
