package response

import (
	"github.com/ubogdan/network-manager-api/model"
)

type Release struct {
	Version string `json:"version"`
}

func FromRelease(rel *model.Release) Release {
	return Release{
		Version: rel.Version.String(),
	}
}
