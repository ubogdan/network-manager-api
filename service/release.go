package service

import (
	"github.com/ubogdan/network-manager-api/repository"
)

// Release service definition.
type Release interface {
	repository.Release
}
