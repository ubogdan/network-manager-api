package request

import (
	models "github.com/ubogdan/network-manager-api/model"
)

type Release struct {
	Version string `json:"version"`
}

func (l *Release) ToModel() models.Release {
	return models.Release{}
}
