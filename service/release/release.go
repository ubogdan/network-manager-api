package release

import (
	models "github.com/ubogdan/network-manager-api/model"
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

func (r *release) Find(channel string) (*models.Release, error) {
	return nil, nil
}

func (r *release) Create(release *models.Release) error {
	return nil
}

func (r *release) Update(release *models.Release) error {
	return nil
}
