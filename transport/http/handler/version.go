package handler

import (
	"net/http"

	"github.com/ubogdan/network-manager-api/service"
	"github.com/ubogdan/network-manager-api/transport/http/response"
)

type version struct {
	log service.Logger
}

func NewVersion(router service.Router, logger service.Logger) {
	handler := version{
		log: logger,
	}
	router.Get("/version", handler.List)

	router.Post("/version", handler.Create)
	router.Delete("/version", handler.Delete)

}

func (h *version) List(w http.ResponseWriter, r *http.Request) error {
	ver := response.Version{
		Version: "0.0.0",
	}
	return response.ToJSON(w, http.StatusOK, ver)
}

//
func (h *version) Create(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *version) Delete(w http.ResponseWriter, r *http.Request) error {
	return nil
}
