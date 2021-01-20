package model

import (
	"errors"
	"time"
)

const (
	ProductName     = "Network Manager Pro"
	DefaultValidity = 30 * 24 * time.Hour
)

var (
	LicenseNotFound = errors.New("License not found")
	LicenseExpired  = errors.New("license ")
)
