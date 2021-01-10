package service

import (
	models "github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/repository"
)

type License interface {
	repository.License
	Renew(license *models.License) error
}
