package response

import (
	"github.com/ubogdan/network-manager-api/model"
)

type License struct {
}

func FromLicese(lic *model.License) License {
	return License{}
}
