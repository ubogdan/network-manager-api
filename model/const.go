package model

import (
	"errors"
)

const (
	ProductName = "Network Manager Pro"
)

var (
	LicenseNotFound = errors.New("License not found")
)
