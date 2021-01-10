package handler

import (
	"net/http"

	models "github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/service"
	"github.com/ubogdan/network-manager-api/transport/http/response"
)

type release struct {
	rel service.Release
	log service.Logger
}

func NewRelease(router service.Router, relSvc service.Release, logger service.Logger) {
	handler := release{
		rel: relSvc,
		log: logger,
	}

	// Release MGMT
	router.Post("/release", handler.Create)
	router.Delete("/release", handler.Delete)

	// Client release channels
	router.Get("/release/latest", handler.List(models.DevelopementChannel))
	router.Get("/release/stable", handler.List(models.ProductionChannel))
}

// List client handler
func (h *release) List(channel string) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		ver := response.Release{
			Version: "0.0.0",
		}
		return response.ToJSON(w, http.StatusOK, ver)
	}
}

func (h *release) Create(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *release) Delete(w http.ResponseWriter, r *http.Request) error {
	return nil
}
