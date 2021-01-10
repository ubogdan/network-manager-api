package service

import (
	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/repository"
)

type License interface {
	repository.License
	Renew(license *model.License) error
}
