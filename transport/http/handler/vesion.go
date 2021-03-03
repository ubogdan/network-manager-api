package handler

import (
	"net/http"

	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/service"
	"github.com/ubogdan/network-manager-api/transport/http/response"
)

func NewVersion(router service.Router) {
	router.Get("/version", Version)
}

func Version(w http.ResponseWriter, _ *http.Request) error {
	return response.ToJSON(w, http.StatusOK, response.Version{Version: model.Version().String()})
}
