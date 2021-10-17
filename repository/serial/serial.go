package serial

import (
	"encoding/pem"
	"errors"
	"reflect"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"

	"github.com/ubogdan/network-manager-api/repository/crypto"
)

var handlerFn func(string, int64) (string, error)

var pemBlock, _ = pem.Decode([]byte(`
-----BEGIN VAULT FILE-----
EP: serial.Generate

mSTwVSsdfeiUGdccOzdufR8y0PfE+BQ3Qz6zbSuk0cAIOKF01hR02vGC/C8CXjJ0
T1HnUDz3B1UxswEKpsr9XSEYYPKlY4WrbN6AS6Ung4LfEQBW9p+bjafJWtQcVCBH
CbVRs37sQ25RLSrgVYetzx9N9Nf7t/QNpumHx+La08VP/Xmdf0qvIGZ9d3NtrDlc
imqXPvExIiE1kzejjtJMalQ6tdIje0o1f0O40cyxlIjPxkjkCQR00tmVwmd2OTKD
8XpFyIgkoliA2seO9ejAkjeGZFWI0E5gW0iUtT3c6U3DvvtM6nliTbsd5HDRNMee
ZOQNLKM9EI1WZd/sD8zx/5Ns17c9CwAKEkh8jIwvnwV1t2AGcrv9JYy/XXHHUTgv
f7iIhIpbOB5k3/CiafQnpP8zWgmQd+YSMzaIXNbMIvXfoV2qsX5xMNNktBd9PtPb
o4o7vLAREdBMQSCYZz/CAbFCrzDygNBabN36rNIh/hdn/kC0hEZFvbJGaTyxHlMd
6Tl11GQLgshRGju+VA5zeNlE5Pw2m4wmlfGYXT5+QxZlRl/ATjzhXu24TztEVW1c
DIclZkwhVqLZ/4kE6PogpJD78BZmZOFwikqgEd3UW8xOQaWkN/wyXE50mAnfcQAb
CT6iKJgZlMCl5kIRpqEAe5NfXV6BL/Zw2HOt3kVTcLIIO8Y3/Dy3/WvTZ1MmVXH3
7jLGyC7XHLh0g2/puugsMnpG3BVtuYFbUOymesAh5EI4ermfnBZOeGlXD4//JE/B
07RCfTtuI+uOgyeQ/UYCmlEJdpn3Zez4viQr7ieoX9pzcj7iowVng5ezMS16GTsE
+9qmhfcx92fUjquZIYLPfYYEeAR8+nesY5xRbfmiEx6JsT6L4imBwSX2J5Ol8TQe
6O6FBCm7ibTUZ7H5LR7QeCEHj/VidVJfHn7xKW8dWDHtPUxDo1Vh+hABNTARk++J
uaMeXRKg//EgghC5Ff51C5e5WojMpzIDNPL4NtcKXjb+RaneblKH+GB0M3TJfz4A
MIRX+Ryc34Q1+swA3Wj9zPcXgzROKDOZTuCDCG+B8od+aHuodAgOMV29SIEakhTe
VTLNsH6Dsl29AWGjbkHAO0nA9K97UzksoxY9SVkJ8EC0Ic6W0uf2b4dRIkNuNqa3
2LYPcAe6j0SYjTTL2HdLV1Tu4PQuJAT8+EhftmXFnyyvI1DgEDIT+r3PtXFqmxol
H2homgOGwScGqEisujXqKmVvkYSqJxM0LE1MLaAIq9eh0KY7lp9DBgIhAYNSLhDM
Vg6/g9iuFlj7d5mJaONgj2BQzcF3KWMnVNnq1Mymo06Rwu12/2wYRpc6Ny8R/4D8
l1uja0iAlqSTwblg/Qq7NnaCsd7MVpDfY/hOrAFklSSRqeLpOv2TH0UoouOfI0fn
cjrNsr1+kcj7zNAxcwxsXLXmSAAK9xVav9LlfbAEqAAmkpDd4lz7E7wg6ZFixf6h
xY8tMaB8hjwy84WsrM5hyXGCBCgL1TccOtH5K/FgjLc0ZWYMF4V6Xb5gQHS9yZGS
/lpdVaXar/UeldmwJv9stdQr316Cj0yaF7OeHqH00UnX3ZZu6s+PEofcxypriuvR
oBhAfmEysBfy3/H4GCeKTSyUdh0NNEqkxZ1v/J1u6FTfjq+k2snn+uX03Jcxg4Sm
l0ExJJdj/9lRePEExhYTpZ4tPGQP/iF2xW4DAvNU1mW/De771mk94LQSHqjsKBQ1
z7yLD+FhYq667P80BUySQ6NqVtd7CvH3
-----END VAULT FILE-----`))

// Generate returns a serial number from hardwareID and validUntil date.
func Generate(privateKey, hardwareID string, validUntilUnix int64) (string, error) {
	if handlerFn != nil {
		return handlerFn(hardwareID, validUntilUnix)
	}

	script, err := crypto.DecryptWithStringKey(privateKey, pemBlock.Bytes)
	if err != nil {
		return "", err
	}

	i := interp.New(interp.Options{})
	err = i.Use(stdlib.Symbols)
	if err != nil {
		return "", err
	}

	_, err = i.Eval(string(script))
	if err != nil {
		return "", err
	}

	callFn, ok := pemBlock.Headers["EP"]
	if !ok {
		callFn = "main.Generate"
	}

	eval, err := i.Eval(callFn)
	if err != nil {
		return "", err
	}

	handlerFn = func(hardwareID string, validUntilUnix int64) (string, error) {
		results := eval.Call([]reflect.Value{
			reflect.ValueOf(privateKey),
			reflect.ValueOf(hardwareID),
			reflect.ValueOf(validUntilUnix),
		})
		if len(results) < 2 {
			return "", errors.New("unexpected response")
		}

		if results[1].Interface() != nil {
			return "", results[1].Interface().(error)
		}

		return results[0].Interface().(string), nil
	}

	return handlerFn(hardwareID, validUntilUnix)
}
