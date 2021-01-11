package model

type License struct {
	ID         uint64 `boltholdKey:"ID"`
	Created    int64
	Expire     int64
	Serial     string
	HardwareID string `boltholdIndex:"HardwareId"`
	Customer   Customer
	Features   []Feature
}
