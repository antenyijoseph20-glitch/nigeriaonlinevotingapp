package handlers

import (
	"net/http"

	"nigeriaonlinevoting/sessions"
)

// Logout destroys the current user session.
func Logout(w http.ResponseWriter, r *http.Request) {

	sessions.DeleteSession(w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
