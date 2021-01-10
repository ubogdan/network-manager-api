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

func New(rel repository.Release) *release {
	return &release{
		Release: rel,
	}
}

func (r *release) Find(channel string) (*model.Release, error) {
	return nil, nil
}

func (r *release) Create(release *model.Release) error {
	return nil
}

func (r *release) Update(release *model.Release) error {
	return nil
}
