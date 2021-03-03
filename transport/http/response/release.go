package response

import (
	"github.com/ubogdan/network-manager-api/model"
)

// Release response DTO.
type Release struct {
	Version string `json:"version"`
}

// FromRelease converts release model to release DTO.
func FromRelease(rel *model.Release) Release {
	return Release{
		Version: rel.Version,
	}
}
