package model

import (
	"errors"
	"time"
)

const (
	// ProductName contains the product name.
	ProductName = "Network Manager Pro"
	// DefaultGracePeriod godoc.
	DefaultGracePeriod = 7 * 24 * time.Hour
	// DefaultValidity godoc.
	DefaultValidity = 30 * 24 * time.Hour

	// ReadLimit1MB used to prevent dos.
	ReadLimit1MB = 1024 * 1024
)

var (
	// ErrLicenseAlreadyExists godoc.
	ErrLicenseAlreadyExists = errors.New("a license for this hardware id already exists")
	// ErrLicenseNotFound godoc.
	ErrLicenseNotFound = errors.New("license not found")
	// ErrLicenseExpired godoc.
	ErrLicenseExpired = errors.New("license expired")
)
