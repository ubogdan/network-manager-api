package model

import (
	"errors"
	"time"
)

const (
	ProductName        = "Network Manager Pro"
	DefaultGracePeriod = 7 * 24 * time.Hour
	DefaultValidity    = 30 * 24 * time.Hour

	ReadLimit1MB  = 1024 * 1024
	ReadLimit10MB = 10 * ReadLimit1MB
)

var (
	ErrLicenseAlreadyExists = errors.New("a license for this hardware id already exists")
	ErrLicenseNotFound      = errors.New("license not found")
	ErrLicenseExpired       = errors.New("license expired")
)
