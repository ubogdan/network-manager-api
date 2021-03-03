package repository

import (
	"github.com/ubogdan/network-manager-api/model"
)

// Release repository definition.
type Release interface {
	Find(channel string) (*model.Release, error)
	Create(release *model.Release) error
	Update(release *model.Release) error
}
