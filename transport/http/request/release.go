package request

import (
	"github.com/ubogdan/network-manager-api/model"
)

type Release struct {
	Version string `json:"version"`
}

func (l *Release) ToModel() model.Release {
	return model.Release{}
}
