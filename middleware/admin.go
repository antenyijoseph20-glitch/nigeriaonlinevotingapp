package middleware

import (
	"net/http"

	"nigeriaonlinevoting/repositories"
	"nigeriaonlinevoting/sessions"
)

// RequireAdmin allows only authenticated administrators.
func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// Get logged-in user ID from the session.
		userID, err := sessions.GetSessionUserID(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Load users from the JSON repository.
		userRepo := repositories.NewJSONRepository("data/users.json")

		user, err := userRepo.GetByID(userID)
		if err != nil {
			http.Error(w, "User not found.", http.StatusUnauthorized)
			return
		}

		// Allow only admins.
		if user.Role != "admin" && user.Role != "super_admin" {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return
		}

		// Continue to the next handler.
		next(w, r)
	}
}
