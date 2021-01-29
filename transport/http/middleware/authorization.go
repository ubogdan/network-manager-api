package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ubogdan/network-manager-api/transport/http/response"
)

func Authorization(authKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.URL.Path, "/v1/renew/") {
				authError := func(authValue string) error {
					authorization := strings.TrimSpace(r.Header.Get("Authorization"))

					if authorization == "" {
						return errors.New("authorization required")
					}

					authparts := strings.Split(authorization, " ")
					if len(authparts) != 2 {
						return errors.New("malformed authorization header")
					}

					switch strings.ToLower(authparts[0]) {
					case "bearer":
						if authparts[1] != authValue {
							return errors.New("malformed authorization header")
						}
						return nil
					default:
						return errors.New(authparts[0] + " not supported ")
					}
				}(authKey)

				if authError != nil {
					response.ToJSON(w, http.StatusUnauthorized, authError.Error())
					return
				}
			}
			next.ServeHTTP(w, r)
		})

	}
}
