package handlers

import (
	"net/http"

	"nigeriaonlinevoting/sessions"
)

// Logout handles user logout.
func Logout(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	sessions.DeleteSession(w, r)

	http.Redirect(
		w,
		r,
		"/login",
		http.StatusSeeOther,
	)
}
