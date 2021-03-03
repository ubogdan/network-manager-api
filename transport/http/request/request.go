package request

import (
	"encoding/json"
	"io"

	"github.com/ubogdan/network-manager-api/model"
)

// FromJSON decode request into specified interface.
func FromJSON(reader io.ReadCloser, payload interface{}) error {
	return UnmarshalWithLimit(reader, model.ReadLimit1MB, payload)
}

// UnmarshalWithLimit decode JSON from request to specified interface using limit reader.
func UnmarshalWithLimit(reader io.ReadCloser, size int64, payload interface{}) error {
	defer reader.Close()

	return json.NewDecoder(io.LimitReader(reader, size)).Decode(&payload)
}
