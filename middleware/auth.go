package middleware

import (
	"net/http"

	"nigeriaonlinevoting/sessions"
)

// RequireAuth protects routes that need login.
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		_, err := sessions.GetSessionUserID(r)

		if err != nil {

			http.Redirect(
				w,
				r,
				"/login",
				http.StatusSeeOther,
			)

			return
		}

		// User is authenticated
		next(w, r)
	}
}
