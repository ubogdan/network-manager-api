package middleware

import (
	"net/http"
	"strings"
)

const (
	corsOrigin  = "Access-Control-Allow-Origin"
	corsHeaders = "Access-Control-Allow-Headers"
	corsMethods = "Access-Control-Allow-Methods"
)

type corsOptions struct {
	Origin  string
	methods []string
	headers []string
}

// corsOptions Middleware
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
			w.Header().Set(corsOrigin, cors.Origin)
			w.Header().Set(corsMethods, allowMethods)
			w.Header().Set(corsHeaders, allowHeaders)
			if r.Method == http.MethodOptions {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func WithMethods(methods ...string) func(*corsOptions) {
	return func(c *corsOptions) {
		c.methods = methods
	}
}

func WithHeaders(headers ...string) func(*corsOptions) {
	return func(c *corsOptions) {
		c.headers = headers
	}
}
