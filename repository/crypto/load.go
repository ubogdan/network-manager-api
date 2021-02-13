package crypto

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

const (
	rsa = "RSA PRIVATE KEY"
	ecc = "EC PRIVATE KEY"
)

// Load rsa or ecc private key from disk
func Load(name string) (crypto.Signer, error) {
	pemBytes, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemBytes)
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
