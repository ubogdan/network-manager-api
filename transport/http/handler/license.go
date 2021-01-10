package handler

import (
	"net/http"

	"github.com/ubogdan/network-manager-api/service"
)

type license struct {
	lic service.License
	log service.Logger
}

func NewLicense(router service.Router, licSvc service.License, logger service.Logger) {
	handler := license{
		lic: licSvc,
		log: logger,
	}

	// License MGMT
	router.Get("/licenses", handler.List)
	router.Post("/licenses", handler.Create)
	router.Put("/licenses/{id}", handler.Update)
	router.Delete("/licenses/{id}", handler.Delete)

	// Client Handler
	router.Put("/license/{serial}/renew", handler.Renew)
}

func (h *license) List(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *license) Create(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *license) Update(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *license) Delete(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Renew client handler
func (h *license) Renew(w http.ResponseWriter, r *http.Request) error {
	return nil
}
