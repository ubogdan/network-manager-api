package request

import (
	"github.com/ubogdan/network-manager-api/model"
)

// Release request DTO.
type Release struct {
	Version string `json:"version"`
}

// ToModel converts release DTO to release model.
func (l *Release) ToModel() model.Release {
	return model.Release{
		Version: l.Version,
	}
}
