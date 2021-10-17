package middleware

import (
	"net/http"
	"strings"
)

const (
	// CorsAllowOriginHeader godoc.
	CorsAllowOriginHeader = "Access-Control-Allow-Origin"
	// CorsAllowHeadersHeader godoc.
	CorsAllowHeadersHeader = "Access-Control-Allow-Headers"
	// CorsAllowMethodsHeader godoc.
	CorsAllowMethodsHeader = "Access-Control-Allow-Methods"
)

type corsOptions struct {
	Origin  string
	methods []string
	headers []string
}

// CORS Middleware.
func CORS(options ...func(*corsOptions)) func(next http.Handler) http.Handler {
	cors := &corsOptions{
		Origin:  "*",
		headers: []string{"Content-Type"},
	}

	for _, optionFn := range options {
		optionFn(cors)
	}

	allowMethods := strings.Join(cors.methods, ", ")

	allowHeaders := strings.Join(cors.headers, ", ")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(CorsAllowOriginHeader, cors.Origin)
			w.Header().Set(CorsAllowMethodsHeader, allowMethods)
			w.Header().Set(CorsAllowHeadersHeader, allowHeaders)
			if r.Method == http.MethodOptions {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// WithMethods godoc.
func WithMethods(methods ...string) func(*corsOptions) {
	return func(c *corsOptions) {
		c.methods = methods
	}
}

// WithHeaders godoc.
func WithHeaders(headers ...string) func(*corsOptions) {
	return func(c *corsOptions) {
		c.headers = headers
	}
}
