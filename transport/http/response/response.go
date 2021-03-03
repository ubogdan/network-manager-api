package response

import (
	"encoding/json"
	"net/http"
)

// ToJSON encode an interface to object.
func ToJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}
