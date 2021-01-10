package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	license2 "github.com/ubogdan/network-manager-api/repository/license"
	release2 "github.com/ubogdan/network-manager-api/repository/release"
	"github.com/ubogdan/network-manager-api/service/license"
	"github.com/ubogdan/network-manager-api/service/release"
	"github.com/ubogdan/network-manager-api/service/router"
	"github.com/ubogdan/network-manager-api/transport/http/handler"
	"github.com/ubogdan/network-manager-api/transport/http/middleware"
)

var listen string
var s3Domain, s3AccessKey, s3SecretKey = os.Getenv("S3_DOMAIN"), os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY")
var licPrivate = os.Getenv("PRIVATE_KEY")

func main() {

	flag.StringVar(&listen, "listen", ":8080", "http listen addres")
	flag.Parse()

	logSvc := logrus.New()

	r := mux.NewRouter()

	api := r.PathPrefix("/v1").Subrouter()
	api.Methods(http.MethodOptions)
	api.Use(
		middleware.CORS(
			middleware.WithHeaders("Content-Type", "Authorization"),
			middleware.WithMethods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete),
		),
		middleware.RateLimit(10), // 10 requests/second
	)

	muxRouter := router.NewMuxRouter(api, logSvc)

	licSvc := license.New(license2.New())
	relSvc := release.New(release2.New())
	handler.NewLicense(muxRouter, licSvc, logSvc)
	handler.NewRelease(muxRouter, relSvc, logSvc)

	// ----------------
	httpd := &http.Server{
		Addr:           listen,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        r,
	}

	httpd.ListenAndServe()
}
