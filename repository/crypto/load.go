package crypto

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

const (
	rsa = "RSA PRIVATE KEY"
	ecc = "EC PRIVATE KEY"
)

// Load rsa or ecc private key from disk.
func Load(key string) (crypto.Signer, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("unable to decode pem data")
	}

	switch block.Type {
	case rsa:
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case ecc:
		return x509.ParseECPrivateKey(block.Bytes)
	default:
		return nil, errors.New("invalid key type")
	}
}
