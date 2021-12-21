package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/ubogdan/network-manager-api/service"
	"github.com/ubogdan/network-manager-api/transport/http/request"
	"github.com/ubogdan/network-manager-api/transport/http/response"
)

// Backup service definition.
type Backup interface {
	GetUploadLink(hardwareID, fileName string, validFor time.Duration) (string, error)
}

type backup struct {
	Backup  Backup
	License License
	Logger  Logger
}

// NewBackup register http endpoints for backup.
func NewBackup(router service.Router, backupSvc Backup, licSvc License, logger Logger) {
	handler := backup{
		Backup:  backupSvc,
		License: licSvc,
		Logger:  logger,
	}

	router.Get("/backups", handler.List)
	router.Post("/backups", handler.Create)

}

// List bakups.
func (h *backup) List(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// Create backup authorization.
func (h *backup) Create(w http.ResponseWriter, req *http.Request) error {
	var reqData request.Backup

	err := request.FromJSON(req.Body, &reqData)
	if err != nil {
		return response.ToJSON(w, http.StatusBadRequest, err)
	}

	if len(reqData.LicenseID) == 0 {
		return response.ToJSON(w, http.StatusNotFound, errors.New("license id is required"))
	}

	if len(reqData.FileName) == 0 {
		return response.ToJSON(w, http.StatusNotFound, errors.New("file name is required"))
	}

	lic, err := h.License.FindBySerial(reqData.LicenseID)
	if err != nil {
		return response.ToJSON(w, http.StatusNotFound, err)
	}

	sign, err := h.Backup.GetUploadLink(lic.HardwareID, reqData.FileName, time.Minute)
	if err != nil {
		return response.ToJSON(w, http.StatusInternalServerError, err)
	}

	return response.ToJSON(w, http.StatusOK, response.Backup{Authorization: sign})
}
