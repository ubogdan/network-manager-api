package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ubogdan/network-manager-api/transport/http/response"
)

// Authorization middleware.
func Authorization(authKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/v1/admin/") {
				authError := func(authValue string) error {
					authorization := strings.TrimSpace(r.Header.Get("Authorization"))

					if authorization == "" {
						return errors.New("authorization required")
					}

					split := strings.Split(authorization, " ")
					if len(split) != 2 {
						return errors.New("malformed authorization header")
					}

					authType := strings.ToLower(split[0])

					switch authType {
					case "bearer":
						if split[1] != authValue {
							return errors.New("malformed authorization header")
						}

						return nil
					default:
						return errors.New(authType + " authorization is not supported ")
					}
				}(authKey)

				if authError != nil {
					_ = response.ToJSON(w, http.StatusUnauthorized, response.Error{
						Code:    http.StatusBadRequest,
						Message: authError.Error(),
					})

					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
