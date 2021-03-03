package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ubogdan/network-manager-api/transport/http/middleware"
)

func TestCORS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		method       string
		allowOrigin  string
		allowHeaders []string
		allowMethods []string
		expected     string
	}{
		{
			method:       http.MethodGet,
			allowOrigin:  "*",
			allowMethods: []string{http.MethodGet},
		},
		{
			method:       http.MethodOptions,
			allowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		},
	}

	for _, test := range tests {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, w.Header().Get(middleware.CorsAllowOriginHeader), test.allowOrigin)
			assert.Equal(t, w.Header().Get(middleware.CorsAllowHeadersHeader), strings.Join(test.allowHeaders, ", "))
			assert.Equal(t, w.Header().Get(middleware.CorsAllowMethodsHeader), strings.Join(test.allowMethods, ", "))
		})

		handlerToTest := middleware.CORS(
			middleware.WithHeaders(test.allowHeaders...),
			middleware.WithMethods(test.allowMethods...),
		)(next)

		// create a mock request to use
		req := httptest.NewRequest(test.method, "http://testing", nil)

		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	}
}
