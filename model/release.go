package model

// Release is the release info.
type Release struct {
	Channel ReleaseChannel `boltholdIndex:"serial"`
	Version string
}

// ReleaseChannel godoc.
type ReleaseChannel string

const (
	// DevelopmentChannel godoc.
	DevelopmentChannel = "development"
	ProductionChannel  = "production"
)
