package model

import (
	"encoding/asn1"
)

// Feature is the license capability.
type Feature struct {
	Name   FeatureType
	Expire int64
	Limit  int64
}

// FeatureType a license feature.
type FeatureType string

const (
	// FeatureTypeMaxNetworks godoc.
	FeatureTypeMaxNetworks = "MaxNetworks"
	FeatureTypeMaxDevices  = "MaxDevices"
	FeatureTypeVPNUsers    = "MaxVPNUsers"
)

// Oid returns the feature oid.
func (f FeatureType) Oid() asn1.ObjectIdentifier {
	switch f {
	case FeatureTypeMaxNetworks:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 3, 6, 1}
	case FeatureTypeMaxDevices:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 3, 6, 2}
	case FeatureTypeVPNUsers:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 3, 6, 11}
	default:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 3, 6}
	}
}

// Description returns a feature description.
func (f FeatureType) Description() string {
	switch f {
	case FeatureTypeMaxNetworks:
		return "Max number of allowed networks can be managed"
	case FeatureTypeMaxDevices:
		return "Max number of allowed devices can be managed"
	case FeatureTypeVPNUsers:
		return "Max number of allowed VPN users"
	default:
		return "Unknown feature"
	}
}
