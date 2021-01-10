package model

type License struct {
	ID         uint64 `boltholdKey:"ID"`
	Serial     string `boltholdIndex:"serial"`
	HardwareID string
	Customer   string
}
