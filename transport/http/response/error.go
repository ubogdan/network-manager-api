package response

import (
	"encoding/json"
	
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(w http.ResponseWriter, status, code int, message string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(Error{
		Code:    code,
		Message: message,
	})
}
