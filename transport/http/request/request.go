package request

import (
	"encoding/json"
	"io"

	"github.com/ubogdan/network-manager-api/model"
)

func FromJSON(reader io.ReadCloser, payload interface{}) error {
	return UnmarshalWithLimit(reader, model.ReadLimit1MB, payload)
}

func UnmarshalWithLimit(reader io.ReadCloser, size int64, payload interface{}) error {
	return json.NewDecoder(io.LimitReader(reader, size)).Decode(&payload)
}
