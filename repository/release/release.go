package release

import (
	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/repository"
)

type release struct {
}

var _ repository.Release = New()

func New() *release {
	return &release{}
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
