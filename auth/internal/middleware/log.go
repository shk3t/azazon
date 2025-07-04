package middleware

import (
	"base/pkg/log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Loggers.Request.Printf(
				"%s %s %s %v",
				r.Method, r.URL.Path, r.RemoteAddr, time.Since(start),
			)
		})
}