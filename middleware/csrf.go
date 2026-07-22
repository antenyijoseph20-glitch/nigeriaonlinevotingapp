package middleware

import "net/http"

// CSRFMiddleware is a placeholder.
// Real CSRF protection will be added later.
func CSRFMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)

	}
}
