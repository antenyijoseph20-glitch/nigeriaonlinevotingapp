package handlers

import (
	"html/template"
	"log"
	"net/http"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/services"
	"nigeriaonlinevoting/sessions"
)

type DashboardHandler struct {
	AuthService         *services.AuthService
	VerificationService *services.VerificationService
}

func NewDashboardHandler(
	authService *services.AuthService,
	verificationService *services.VerificationService,
) *DashboardHandler {

	return &DashboardHandler{
		AuthService:         authService,
		VerificationService: verificationService,
	}
}

func (h *DashboardHandler) Dashboard(
	w http.ResponseWriter,
	r *http.Request,
) {

	// Get logged-in user ID
	userID, err := sessions.GetSessionUserID(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Load user
	user, err := h.AuthService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found.", http.StatusNotFound)
		return
	}

	// Default verification status
	verificationStatus := "not_submitted"

	verification, err := h.VerificationService.GetUserVerification(userID)
	if err == nil && verification != nil {
		verificationStatus = verification.Status
	}

	// Data sent to dashboard.html
	data := struct {
		User               *models.User
		VerificationStatus string
	}{
		User:               user,
		VerificationStatus: verificationStatus,
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, "Unable to load dashboard.", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println("dashboard template execution:", err)
		return
	}
}
