package response

import (
	models "github.com/ubogdan/network-manager-api/model"
)

type Release struct {
	Version string `json:"version"`
}

func FromRelease(rel *models.Release) Release {
	return Release{
		Version: rel.Version.String(),
	}
}
