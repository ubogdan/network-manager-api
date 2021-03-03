package handler

import (
	"net/http"

	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/service"
	"github.com/ubogdan/network-manager-api/transport/http/request"
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
	router.Put("/release", handler.Update)

	// Client release channels
	router.Get("/release/latest", handler.List(model.DevelopmentChannel))
	router.Get("/release/stable", handler.List(model.ProductionChannel))
}

// List releases.
func (h *release) List(channel string) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		release, err := h.rel.Find(channel)
		if err != nil {
			return response.ToJSON(w, http.StatusInternalServerError, nil)
		}

		return response.ToJSON(w, http.StatusOK, response.FromRelease(release))
	}
}

// Create release.
func (h *release) Create(w http.ResponseWriter, r *http.Request) error {
	var release request.Release

	err := request.FromJSON(r.Body, &release)
	if err != nil {
		return response.ToJSON(w, http.StatusBadRequest, err)
	}

	toModel := release.ToModel()

	err = h.rel.Create(&toModel)
	if err != nil {
		return response.ToJSON(w, http.StatusInternalServerError, err)
	}

	return response.ToJSON(w, http.StatusOK, response.FromRelease(&toModel))
}

func (h *release) Update(w http.ResponseWriter, r *http.Request) error {
	return nil
}
