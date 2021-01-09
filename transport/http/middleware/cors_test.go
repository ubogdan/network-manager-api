package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
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
			assert.Equal(t, w.Header().Get(corsOrigin), test.allowOrigin)
			assert.Equal(t, w.Header().Get(corsHeaders), strings.Join(test.allowHeaders, ", "))
			assert.Equal(t, w.Header().Get(corsMethods), strings.Join(test.allowMethods, ", "))
		})

		handlerToTest := CORS(
			WithHeaders(test.allowHeaders...),
			WithMethods(test.allowMethods...),
		)(next)

		// create a mock request to use
		req := httptest.NewRequest(test.method, "http://testing", nil)

		// call the handler using a mock response recorder (we'll not use that anyway)
		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)

	}
}
