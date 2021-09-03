package main

import (
	"context"
	"encoding/hex"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/ubogdan/network-manager-api/repository/crypto"
	"github.com/ubogdan/network-manager-api/repository/dynamo"
	"github.com/ubogdan/network-manager-api/service/license"
	"github.com/ubogdan/network-manager-api/service/router"
	"github.com/ubogdan/network-manager-api/transport/http/handler"
	"github.com/ubogdan/network-manager-api/transport/http/middleware"
)

func main() {
	licenseKey, authorizedKey, authorizationKey := os.Getenv("LICENSE_KEY"), os.Getenv("AUTHORIZED_KEY"), os.Getenv("API_BEARER_AUTH")

	logSvc := logrus.New()

	authorizedKeyBytes, err := hex.DecodeString(authorizedKey)
	if err != nil {
		logSvc.Fatalf("invalid key %s", err)
	}

	r := mux.NewRouter()
	api := r.PathPrefix("/" + os.Getenv("API_BASE_PATH")).Subrouter()
	api.Methods(http.MethodOptions)
	api.Use(
		middleware.CORS(
			middleware.WithHeaders("Content-Type", "Authorization"),
			middleware.WithMethods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete),
		),
		middleware.Authorization(authorizationKey),
	)

	licenseSigner, err := crypto.Load(licenseKey)
	if err != nil {
		logSvc.Fatalf("read license key %s", err)
	}

	sess, err := session.NewSession()
	if err != nil {
		logSvc.Fatalf("new session %s", err)

	}
	db := dynamodb.New(sess)
	licSvc := license.New(dynamo.License(db), authorizedKeyBytes, licenseSigner)
	//relSvc := release.New(bolt.Release())

	muxRouter := router.NewMuxRouter(api, logSvc)
	handler.NewLicense(muxRouter, licSvc, authorizedKeyBytes, logSvc)
	//handler.NewRelease(muxRouter, relSvc, logSvc)
	handler.NewVersion(muxRouter)

	handler := gorillamux.New(r)
	lambda.Start(
		func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			log.Printf("%s %s", request.HTTPMethod, request.Path)
			return handler.ProxyWithContext(ctx, request)
		},
	)
}
