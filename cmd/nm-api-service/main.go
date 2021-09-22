package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/repository/crypto"
	"github.com/ubogdan/network-manager-api/repository/dynamo"
	"github.com/ubogdan/network-manager-api/service"
	"github.com/ubogdan/network-manager-api/service/license"
	"github.com/ubogdan/network-manager-api/service/router"
	"github.com/ubogdan/network-manager-api/transport/http/handler"
	"github.com/ubogdan/network-manager-api/transport/http/middleware"
)

func main() {
	logSvc := logrus.New()

	_ = xray.Configure(xray.Config{
		DaemonAddr:     os.Getenv("AWS_XRAY_DAEMON_ADDRESS"),
		ServiceVersion: model.Version().String(),
	})

	// Trace init start-up time
	_, startSeg := xray.BeginSegment(context.Background(), "start-up")
	r, err := setup(logSvc)
	startSeg.Close(err)

	// Log to stderr if something goes wrong during initialization
	if err != nil {
		logSvc.Fatalf("start-up error %s", err)

		return
	}

	// Proxy handler
	proxyRouter := gorillamux.New(r)

	// Start lambda event handler
	lambda.Start(
		func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			_, seg := xray.BeginSegment(ctx, request.HTTPMethod+" "+request.Path)

			response, err := proxyRouter.ProxyWithContext(ctx, request)

			seg.HTTP = &xray.HTTPData{
				Request: &xray.RequestData{
					Method: request.HTTPMethod,
					URL:    request.Path,
					Traced: true,
				},
				Response: &xray.ResponseData{
					Status:        response.StatusCode,
					ContentLength: len(response.Body),
				},
			}
			seg.Close(err)

			return response, err
		},
	)
}

func setup(logger service.Logger) (*mux.Router, error) {
	authorizedKeyBytes, err := hex.DecodeString(os.Getenv("AUTHORIZED_KEY"))
	if err != nil {
		return nil, fmt.Errorf("invalid key %s", err)
	}

	r := mux.NewRouter()
	api := r.PathPrefix("/" + os.Getenv("API_BASE_PATH")).Subrouter()
	api.Methods(http.MethodOptions)
	api.Use(
		middleware.CORS(
			middleware.WithHeaders("Content-Type", "Authorization"),
			middleware.WithMethods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete),
		),
		middleware.Authorization(os.Getenv("API_BEARER_AUTH")),
	)

	licenseSigner, err := crypto.Load(os.Getenv("LICENSE_KEY"))
	if err != nil {
		return nil, fmt.Errorf("read license key %s", err)
	}

	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("new session %s", err)
	}

	db := dynamodb.New(sess)

	licSvc := license.New(dynamo.License(db), authorizedKeyBytes, licenseSigner)
	//relSvc := release.New(bolt.Release())

	muxRouter := router.NewMuxRouter(api, logger)
	handler.NewLicense(muxRouter, licSvc, authorizedKeyBytes, logger)
	//handler.NewRelease(muxRouter, relSvc, logger)
	handler.NewVersion(muxRouter)

	return r, nil
}
