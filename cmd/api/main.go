package main

import (
	"encoding/hex"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/ubogdan/network-manager-api/pkg/bolthold"
	"github.com/ubogdan/network-manager-api/repository/bolt"
	"github.com/ubogdan/network-manager-api/repository/crypto"
	"github.com/ubogdan/network-manager-api/service/license"
	"github.com/ubogdan/network-manager-api/service/release"
	"github.com/ubogdan/network-manager-api/service/router"
	"github.com/ubogdan/network-manager-api/transport/http/handler"
	"github.com/ubogdan/network-manager-api/transport/http/middleware"
)

func main() {
	var listen, licenseKey, boltdb string

	flag.StringVar(&boltdb, "database", "netmgrapi.db", "")
	flag.StringVar(&licenseKey, "sign", "signing.key", "")
	flag.StringVar(&listen, "listen", ":8080", "http listen addres")
	flag.Parse()

	logSvc := logrus.New()

	s3AccessKey, s3SecretKey := os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY")
	authorizedKey, authorizationKey := os.Getenv("AUTHORIZED_KEY"), os.Getenv("API_BEARER_AUTH")

	authorizedKeyBytes, err := hex.DecodeString(authorizedKey)
	if err != nil {
		logSvc.Fatalf("invalid key %s", err)
	}

	licenseSigner, err := crypto.Load(licenseKey)
	if err != nil {
		logSvc.Fatalf("read license key %s", err)
	}

	var db *bolthold.Store
	if s3AccessKey != "" && s3SecretKey != "" {
		awsSession, e := session.NewSession(&aws.Config{
			Credentials:      credentials.NewStaticCredentials(s3AccessKey, s3SecretKey, ""),
			Region:           aws.String(os.Getenv("S3_REGION")),
			Endpoint:         aws.String(os.Getenv("S3_ENDPOINT")),
			S3ForcePathStyle: aws.Bool(true),
		})
		if e != nil {
			logSvc.Fatalf("session.NewSession %s", e)
		}
		s3boltdb := bolthold.WithS3(s3.New(awsSession), os.Getenv("S3_BUCKET"), os.Getenv("S3_PREFIX"))
		s3boltdb.Log = logSvc
		db, err = s3boltdb.Open(boltdb, 0755, nil)
	} else {
		logSvc.Infof("warning: s3 sync disabled.")
		db, err = bolthold.Open(boltdb, 0755, nil)
	}
	if err != nil {
		logSvc.Fatalf("bolthold.Open %s", err)
	}
	defer db.Close()

	r := mux.NewRouter()
	api := r.PathPrefix("/v1").Subrouter()
	api.Methods(http.MethodOptions)
	api.Use(
		middleware.CORS(
			middleware.WithHeaders("Content-Type", "Authorization"),
			middleware.WithMethods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete),
		),
		middleware.RateLimit(1), // 1 requests/second
		middleware.Authorization(authorizationKey),
	)

	licSvc := license.New(bolt.License(db), authorizedKeyBytes, licenseSigner)
	relSvc := release.New(bolt.Release())

	muxRouter := router.NewMuxRouter(api, logSvc)
	handler.NewLicense(muxRouter, licSvc, authorizedKeyBytes, logSvc)
	handler.NewRelease(muxRouter, relSvc, logSvc)
	handler.NewVersion(muxRouter)

	// ----------------
	httpd := &http.Server{
		Addr:           listen,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        r,
	}

	err = httpd.ListenAndServe()
	if err != nil {
		logSvc.Errorf("httpd.Listen %s", err)
	}
}
