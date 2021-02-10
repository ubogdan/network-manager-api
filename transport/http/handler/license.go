package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/ubogdan/network-manager-api/service"
	"github.com/ubogdan/network-manager-api/transport/http/request"
	"github.com/ubogdan/network-manager-api/transport/http/response"
)

type license struct {
	secretKey []byte
	license   service.License
	log       service.Logger
}

// TODO replace secretKey with secret manager
func NewLicense(router service.Router, licSvc service.License, secretKey []byte, logger service.Logger) {
	handler := license{
		license:   licSvc,
		log:       logger,
		secretKey: secretKey,
	}

	// LicenseRenew MGMT
	router.Get("/licenses", handler.List)
	router.Get("/licenses/{id}", handler.Find)
	router.Post("/licenses", handler.Create)
	router.Put("/licenses/{id}", handler.Update)
	router.Delete("/licenses/{id}", handler.Delete)

	// Client Handler
	router.Put("/renew/{serial}", handler.Renew)
}

func (h *license) List(w http.ResponseWriter, r *http.Request) error {
	list, err := h.license.FindAll()
	if err != nil {
		return response.ToJSON(w, http.StatusInternalServerError, nil)
	}
	return response.ToJSON(w, http.StatusOK, list)
}

func (h *license) Create(w http.ResponseWriter, r *http.Request) error {
	var lic request.License
	err := request.FromJSON(r.Body, &lic)
	if err != nil {
		return response.ToJSON(w, http.StatusBadRequest, err)
	}
	model := lic.ToModel()
	err = h.license.Create(&model)
	if err != nil {
		return response.ToJSON(w, http.StatusInternalServerError, err)
	}
	return response.ToJSON(w, http.StatusOK, response.FromLicese(&model))
}

func (h *license) Find(w http.ResponseWriter, r *http.Request) error {
	licenseID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		return response.ToJSON(w, http.StatusBadRequest, "invalid license id")
	}

	lic, err := h.license.Find(licenseID)
	if err != nil {
		return response.ToJSON(w, http.StatusInternalServerError, err)
	}

	return response.ToJSON(w, http.StatusOK, response.FromLicese(lic))

}

func (h *license) Update(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	licenseID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return response.ToJSON(w, http.StatusBadRequest, "invalid license id")
	}

	var lic request.License
	err = request.FromJSON(r.Body, &lic)
	if err != nil {
		return response.ToJSON(w, http.StatusBadRequest, err)
	}
	model := lic.ToModel()
	model.ID = licenseID

	err = h.license.Update(&model)
	if err != nil {
		return response.ToJSON(w, http.StatusInternalServerError, err)
	}
	return response.ToJSON(w, http.StatusOK, response.FromLicese(&model))
}

func (h *license) Delete(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	licenseID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return response.ToJSON(w, http.StatusBadRequest, "invalid license id")
	}

	err = h.license.Delete(licenseID)
	if err != nil {
		return response.ToJSON(w, http.StatusInternalServerError, err)
	}
	return response.ToJSON(w, http.StatusOK, "")
}

// Renew client handler
func (h *license) Renew(w http.ResponseWriter, r *http.Request) error {

	serial := strings.TrimSpace(mux.Vars(r)["serial"])
	if serial == "" {
		return response.ToJSON(w, http.StatusBadRequest, "invalid serial number")
	}

	renew, err := request.LicenseFromEncryptedPayload(r.Body, h.secretKey)
	if err != nil {
		return response.ToJSON(w, http.StatusInternalServerError, "invalid license data")
	}

	if serial != renew.Serial {
		return response.ToJSON(w, http.StatusInternalServerError, "serial number mismatch")
	}

	//
	model := renew.ToModel()

	derBytes, err := h.license.Renew(&model)
	if err != nil {
		return response.ToJSON(w, http.StatusInternalServerError, "renewal denied")
	}

	return response.LicenseToEncryptedPayload(w, derBytes, h.secretKey)
}
