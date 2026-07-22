package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs every incoming request.
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		log.Printf(
			"%s %s",
			r.Method,
			r.URL.Path,
		)

		next.ServeHTTP(w, r)

		log.Printf(
			"Completed in %v",
			time.Since(start),
		)
	}
}
