package release

import (
	models "github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/repository"
)

type release struct {
}

var _ repository.Release = New()

func New() *release {
	return &release{}
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
