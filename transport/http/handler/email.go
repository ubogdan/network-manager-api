package handler

import (
	"net/http"
	"net/mail"
	"os"
	"time"

	"github.com/ubogdan/network-manager-api/service"
	"github.com/ubogdan/network-manager-api/service/email"
	"github.com/ubogdan/network-manager-api/transport/http/request"
	"github.com/ubogdan/network-manager-api/transport/http/response"
)

func NewEmail(router service.Router, licSvc License, logger Logger) {
	// Email notification
	router.Post("/notify/email", handleEmailNotifcation(licSvc, logger))
}

func handleEmailNotifcation(license License, logger Logger) func(http.ResponseWriter, *http.Request) error {
	from := mail.Address{
		Name:    "Self Service",
		Address: os.Getenv("EMAIL_FROM"),
	}
	if from.Address == "" {
		logger.Errorf("missing EMAIL_FROM environment variable")
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		var req request.Email

		err := request.FromJSON(r.Body, &req)
		if err != nil {
			logger.Errorf("invalid payload %s", err)
			return response.NewError(w, http.StatusBadRequest, 1040, "invalid payload")
		}

		switch req.Product {
		case "47fe5af0-f722-11e9-8f0b-362b9e155667": // Network Manager Pro
			license, err := license.FindBySerial(req.License)
			if err != nil {
				return response.NewError(w, http.StatusUnauthorized, 1041, "invalid license")
			}
			// Validate license
			if time.Unix(license.Expire, 0).After(time.Now()) {
				logger.Errorf("expired license expire= %d", license.Expire)
				return response.NewError(w, http.StatusUnauthorized, 1042, "license expired")
			}

			// Validate to address
			to, err := mail.ParseAddress(req.To)
			if err != nil {
				logger.Errorf("failed to parse email address: %v", err)
				return response.NewError(w, http.StatusBadRequest, 1043, "invalid email address")
			}

			err = email.Notify(from.String(), to.String(), req.Subject, req.ToModel(), email.Template{
				Product: email.Product{
					Name: "Network Manager Pro",
					Link: "https://lfpanels.com/",
					Logo: os.Getenv("EMAIL_LOGO"),
				},
			})
			if err != nil {
				logger.Errorf("failed to send email: %v", err)
				return response.NewError(w, http.StatusInternalServerError, 1044, "failed to send email")
			}

		default:
			return response.NewError(w, http.StatusUnauthorized, 1050, "invalid or missing product id")
		}

		return response.ToJSON(w, http.StatusOK, nil)
	}
}
