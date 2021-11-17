package handler

import (
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/service"

	"github.com/ubogdan/network-manager-api/transport/http/request"
	"github.com/ubogdan/network-manager-api/transport/http/response"
)

// License service definition.
type License interface {
	FindAll() ([]model.License, error)
	Find(id string) (*model.License, error)
	Create(license *model.License) error
	Update(license *model.License) error
	Delete(id string) error
	Renew(license *model.License) ([]byte, error)
}

// Logger service definition.
type Logger interface {
	Errorf(string, ...interface{})
}

type license struct {
	secretKey []byte
	license   License
	log       Logger
}

// NewLicense register http endpoints for license.
func NewLicense(router service.Router, licSvc License, secretKey []byte, logger Logger) {
	handler := license{
		license:   licSvc,
		log:       logger,
		secretKey: secretKey,
	}

	// Renew MGMT
	router.Get("/licenses", handler.List)
	router.Get("/licenses/{id}", handler.Find)
	router.Post("/licenses", handler.Create)
	router.Put("/licenses/{id}", handler.Update)
	router.Delete("/licenses/{id}", handler.Delete)

	// Client Handler
	router.Put("/renew/{serial}", handler.Renew)
}

// List returns a list of licenses.
func (h *license) List(w http.ResponseWriter, _ *http.Request) error {
	list, err := h.license.FindAll()
	if err != nil {
		return response.ToJSON(w, http.StatusInternalServerError, nil)
	}

	return response.ToJSON(w, http.StatusOK, response.FromLicenses(list))
}

// Create register a new license.
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

	return response.ToJSON(w, http.StatusOK, response.FromLicense(&model))
}

// Find license by id.
func (h *license) Find(w http.ResponseWriter, r *http.Request) error {
	licenseID := mux.Vars(r)["id"]
	_, err := hex.DecodeString(licenseID)
	if err != nil {
		return response.ToJSON(w, http.StatusBadRequest, "invalid license id")
	}

	lic, err := h.license.Find(licenseID)
	if err != nil {
		return response.ToJSON(w, http.StatusInternalServerError, err)
	}

	return response.ToJSON(w, http.StatusOK, response.FromLicense(lic))
}

// Update license details.
func (h *license) Update(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	licenseID := params["id"]
	_, err := hex.DecodeString(licenseID)
	if err != nil {
		return response.ToJSON(w, http.StatusBadRequest, "invalid license id")
	}

	var lic request.License

	err = request.FromJSON(r.Body, &lic)
	if err != nil {
		return response.ToJSON(w, http.StatusBadRequest, err)
	}

	model := lic.ToModel()

	err = h.license.Update(&model)
	if err != nil {
		return response.ToJSON(w, http.StatusInternalServerError, err)
	}

	return response.ToJSON(w, http.StatusOK, response.FromLicense(&model))
}

// Delete license.
func (h *license) Delete(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	licenseID := params["id"]

	_, err := hex.DecodeString(licenseID)
	if err != nil {
		return response.ToJSON(w, http.StatusBadRequest, "invalid license id")
	}

	err = h.license.Delete(licenseID)
	if err != nil {
		return response.ToJSON(w, http.StatusInternalServerError, err)
	}

	return response.ToJSON(w, http.StatusOK, "")
}

// Renew client handler.
func (h *license) Renew(w http.ResponseWriter, r *http.Request) error {
	serial := strings.TrimSpace(mux.Vars(r)["serial"])
	if serial == "" {
		return response.ToJSON(w, http.StatusBadRequest, "invalid serial number")
	}

	renew, err := request.LicenseFromEncryptedPayload(r.Body, h.secretKey)
	if err != nil {
		return response.ToJSON(w, http.StatusBadRequest, "invalid license data")
	}

	if serial != renew.Serial {
		return response.ToJSON(w, http.StatusBadRequest, "serial number mismatch")
	}

	toModel := renew.ToModel()

	license, err := h.license.Find(toModel.HardwareID)
	if err != nil {
		return response.ToJSON(w, http.StatusNotFound, "no license match your hardware id")
	}

	if toModel.Serial != license.Serial {
		return response.ToJSON(w, http.StatusNotFound, "unexpected serial number")
	}

	derBytes, err := h.license.Renew(license)
	if err != nil {
		h.log.Errorf("license.Renew %s: serial:%s hardwareID:%s error:%s",
			r.RemoteAddr, toModel.Serial, toModel.HardwareID, err)

		return response.ToJSON(w, http.StatusInternalServerError, "renewal denied")
	}

	return response.LicenseToEncryptedPayload(w, derBytes, h.secretKey)
}
