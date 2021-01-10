package handler

import (
	"net/http"

	models "github.com/ubogdan/network-manager-api/model"
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
	router.Get("/release/latest", handler.List(models.DevelopementChannel))
	router.Get("/release/stable", handler.List(models.ProductionChannel))
}

// List client handler
func (h *release) List(channel string) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		release, err := h.rel.Find(channel)
		if err != nil {
			return response.ToJSON(w, http.StatusInternalServerError, nil)
		}
		return response.ToJSON(w, http.StatusOK, response.FromRelease(release))
	}
}

func (h *release) Create(w http.ResponseWriter, r *http.Request) error {
	var rel request.Release
	err := request.FromJSON(r.Body, &rel)
	if err != nil {
		return response.ToJSON(w, http.StatusBadRequest, err)
	}
	model := rel.ToModel()
	err = h.rel.Create(&model)
	if err != nil {
		return response.ToJSON(w, http.StatusInternalServerError, err)
	}
	return response.ToJSON(w, http.StatusOK, response.FromRelease(&model))
}

func (h *release) Update(w http.ResponseWriter, r *http.Request) error {
	return nil
}
