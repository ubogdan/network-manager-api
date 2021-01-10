package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/ubogdan/network-manager-api/service/router"
	"github.com/ubogdan/network-manager-api/transport/http/handler"
	"github.com/ubogdan/network-manager-api/transport/http/middleware"
)

var listen string
var s3Domain, s3AccessKey, s3SecretKey = os.Getenv("S3_DOMAIN"), os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY")

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"version":"0.3.1"}`))
}

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

	handler.NewLicense(muxRouter, nil, logSvc)
	handler.NewRelease(muxRouter, nil, logSvc)

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
