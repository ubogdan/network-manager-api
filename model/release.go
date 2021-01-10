package models

import (
	"github.com/coreos/go-semver/semver"
)

type Release struct {
	Version semver.Version
}

type ReleaseChannel string

const (
	DevelopementChannel = "development"
	ProductionChannel   = "production"
)
