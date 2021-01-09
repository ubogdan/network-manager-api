package request

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func FromJSON(reader io.ReadCloser, payload interface{}) error {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	defer reader.Close()

	return json.Unmarshal(body, &payload)
}
