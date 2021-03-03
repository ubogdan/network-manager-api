package release

import (
	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/repository"
	"github.com/ubogdan/network-manager-api/service"
)

type release struct {
	Release repository.Release
}

var _ service.Release = New(nil)

// New godoc.
func New(rel repository.Release) *release {
	return &release{
		Release: rel,
	}
}

// Find godoc.
func (r *release) Find(channel string) (*model.Release, error) {
	return nil, nil
}

// Create godoc.
func (r *release) Create(release *model.Release) error {
	return nil
}

// Update godoc.
func (r *release) Update(release *model.Release) error {
	return nil
}
