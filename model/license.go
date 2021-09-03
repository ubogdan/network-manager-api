package model

// License is the license model.
type License struct {
	Created    int64
	Expire     int64
	LastIssued int64
	Serial     string
	HardwareID string
	Customer   Customer
	Features   []Feature
}
