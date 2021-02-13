package bolt

import (
	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/repository"
)

type release struct {
}

var _ repository.Release = Release()

// Release return a release repository
func Release() *release {
	return &release{}
}

// Find godoc
func (r *release) Find(channel string) (*model.Release, error) {
	return nil, nil
}

// Create godoc
func (r *release) Create(release *model.Release) error {
	return nil
}

// Update godoc
func (r *release) Update(release *model.Release) error {
	return nil
}
