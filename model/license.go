package model

// License is the license model
type License struct {
	ID         uint64 `boltholdKey:"ID"`
	Created    int64
	Expire     int64
	LastIssued int64
	Serial     string
	HardwareID string `boltholdIndex:"HardwareId"`
	Customer   Customer
	Features   []Feature
}
