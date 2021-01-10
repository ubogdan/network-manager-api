package response

import (
	models "github.com/ubogdan/network-manager-api/model"
)

type License struct {
}

func FromLicese(lic *models.License) License {
	return License{}
}
