package middleware

import "net/http"

// RateLimitMiddleware is a placeholder.
// We'll implement IP-based request limiting later.
func RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)

	}
}
