package handlers

import (
	"html/template"
	"net/http"

	"nigeriaonlinevoting/services"
	"nigeriaonlinevoting/sessions"
)

type DashboardHandler struct {
	AuthService *services.AuthService
}

func NewDashboardHandler(service *services.AuthService) *DashboardHandler {
	return &DashboardHandler{
		AuthService: service,
	}
}

func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {

	// Read session cookie
	userID, err := sessions.GetSessionUserID(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Load user information
	user, err := h.AuthService.GetUserByID(userID)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Load dashboard template
	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, "Unable to load dashboard", http.StatusInternalServerError)
		return
	}

	// Display logged-in user
	if err := tmpl.Execute(w, user); err != nil {
		http.Error(w, "Unable to render dashboard", http.StatusInternalServerError)
	}
}
