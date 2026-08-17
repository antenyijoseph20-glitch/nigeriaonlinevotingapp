package handlers

import (
	"html/template"
	"net/http"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/services"
	"nigeriaonlinevoting/sessions"
)

type VerificationHandler struct {
	VerificationService *services.VerificationService
}

func NewVerificationHandler(
	service *services.VerificationService,
) *VerificationHandler {

	return &VerificationHandler{
		VerificationService: service,
	}
}

// Verification handles voter verification.
func (h *VerificationHandler) Verification(
	w http.ResponseWriter,
	r *http.Request,
) {

	// Ensure the user is logged in.
	userID, err := sessions.GetSessionUserID(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	switch r.Method {

	case http.MethodGet:

		// Get the logged-in user.
		user, err := h.VerificationService.GetUserByID(userID)
		if err != nil {
			http.Error(w, "User not found.", http.StatusNotFound)
			return
		}

		// Data sent to the template.
		data := struct {
			FirstName string
			LastName  string
			Email     string
			NIN       string
		}{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			NIN:       user.NIN,
		}

		tmpl, err := template.ParseFiles("templates/verification.html")
		if err != nil {
			http.Error(w, "Unable to load verification page.", http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Unable to render page.", http.StatusInternalServerError)
			return
		}

	case http.MethodPost:

		fullName := r.FormValue("fullName")
		dateOfBirth := r.FormValue("dateOfBirth")
		gender := r.FormValue("gender")
		state := r.FormValue("state")
		lga := r.FormValue("lga")
		ward := r.FormValue("ward")
		pollingUnit := r.FormValue("pollingUnit")

		// Validation
		if fullName == "" {
			http.Error(w, "Full name is required.", http.StatusBadRequest)
			return
		}

		if dateOfBirth == "" {
			http.Error(w, "Date of birth is required.", http.StatusBadRequest)
			return
		}

		if gender == "" {
			http.Error(w, "Gender is required.", http.StatusBadRequest)
			return
		}

		if state == "" {
			http.Error(w, "State is required.", http.StatusBadRequest)
			return
		}

		if lga == "" {
			http.Error(w, "LGA is required.", http.StatusBadRequest)
			return
		}

		if ward == "" {
			http.Error(w, "Ward is required.", http.StatusBadRequest)
			return
		}

		if pollingUnit == "" {
			http.Error(w, "Polling Unit is required.", http.StatusBadRequest)
			return
		}

		verification := models.Verification{
			UserID:      userID,
			FullName:    fullName,
			DateOfBirth: dateOfBirth,
			Gender:      gender,
			State:       state,
			LGA:         lga,
			Ward:        ward,
			PollingUnit: pollingUnit,
		}

		err = h.VerificationService.SubmitVerification(verification)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Return to dashboard after successful submission.
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
