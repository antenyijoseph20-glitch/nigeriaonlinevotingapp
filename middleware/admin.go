package middleware

import "net/http"

// AdminMiddleware protects administrator routes.
func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// Later we'll verify:
		// 1. User is logged in.
		// 2. User role is "admin".

		next.ServeHTTP(w, r)

	}
}
