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
	FeatureTypeMaxNetworks = "MaxNetworks"
)

// Oid returns the feature oid.
func (f FeatureType) Oid() asn1.ObjectIdentifier {
	switch f {
	case FeatureTypeMaxNetworks:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 3, 6, 1}
	default:
		return asn1.ObjectIdentifier{1, 3, 6, 1, 3, 6}
	}
}

// Description returns a feature description.
func (f FeatureType) Description() string {
	switch f {
	case FeatureTypeMaxNetworks:
		return "Max number of allowed networks can be managed"
	default:
		return ""
	}
}
