package handlers

import (
	"html/template"
	"net/http"

	"nigeriaonlinevoting/services"
	"nigeriaonlinevoting/sessions"
)

type ProfileHandler struct {
	AuthService *services.AuthService
}

func NewProfileHandler(service *services.AuthService) *ProfileHandler {
	return &ProfileHandler{
		AuthService: service,
	}
}

// Profile displays and updates the voter's profile.
func (h *ProfileHandler) Profile(w http.ResponseWriter, r *http.Request) {

	// Get the logged-in user's ID.
	userID, err := sessions.GetSessionUserID(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Load the user.
	user, err := h.AuthService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found.", http.StatusNotFound)
		return
	}

	switch r.Method {

	case http.MethodGet:

		tmpl, err := template.ParseFiles("templates/profile.html")
		if err != nil {
			http.Error(w, "Unable to load profile page.", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, user)

	case http.MethodPost:

		user.OtherName = r.FormValue("otherName")
		user.Gender = r.FormValue("gender")
		user.State = r.FormValue("state")
		user.LGA = r.FormValue("lga")
		user.Ward = r.FormValue("ward")
		user.PollingUnit = r.FormValue("pollingUnit")
		user.VIN = r.FormValue("vin")

		if err := h.AuthService.UpdateUser(*user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
