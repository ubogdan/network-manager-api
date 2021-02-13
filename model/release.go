package model

// Release is the release info
type Release struct {
	Channel ReleaseChannel `boltholdIndex:"serial"`
	Version string
}

type ReleaseChannel string

const (
	DevelopementChannel = "development"
	ProductionChannel   = "production"
)
