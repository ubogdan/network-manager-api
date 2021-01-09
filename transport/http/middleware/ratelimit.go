package middleware

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"go.uber.org/ratelimit"
)

func RateLimit(rate int, opts ...ratelimit.Option) func(next http.Handler) http.Handler {
	var lmap sync.Map

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, fmt.Sprintf("invalid RemoteAddr: %s", err), http.StatusInternalServerError)
				return
			}
			lif, ok := lmap.Load(host)
			if !ok {
				lif = ratelimit.New(rate, opts...)
			}
			lm := lif.(ratelimit.Limiter)
			lmap.Store(host, lm)
			lm.Take()
			next.ServeHTTP(w, r)
		})
	}
}
