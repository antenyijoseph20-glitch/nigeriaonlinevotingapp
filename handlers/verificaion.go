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

// Verification page handler
func (h *VerificationHandler) Verification(
	w http.ResponseWriter,
	r *http.Request,
) {

	switch r.Method {

	case http.MethodGet:

		tmpl, err := template.ParseFiles(
			"templates/verification.html",
		)

		if err != nil {

			http.Error(
				w,
				"Unable to load verification page",
				http.StatusInternalServerError,
			)

			return
		}

		tmpl.Execute(
			w,
			nil,
		)

	case http.MethodPost:

		// Get logged-in user ID
		userID, err := sessions.GetSessionUserID(r)

		if err != nil {

			http.Redirect(
				w,
				r,
				"/login",
				http.StatusSeeOther,
			)

			return
		}

		// Read form data

		fullName := r.FormValue("fullName")
		dateOfBirth := r.FormValue("dateOfBirth")
		gender := r.FormValue("gender")

		state := r.FormValue("state")
		lga := r.FormValue("lga")
		ward := r.FormValue("ward")
		pollingUnit := r.FormValue("pollingUnit")

		verification := models.Verification{

			UserID: userID,

			FullName: fullName,

			DateOfBirth: dateOfBirth,

			Gender: gender,

			State: state,

			LGA: lga,

			Ward: ward,

			PollingUnit: pollingUnit,
		}

		err = h.VerificationService.SubmitVerification(
			verification,
		)

		if err != nil {

			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)

			return
		}

		tmpl, err := template.ParseFiles(
			"templates/verification_success.html",
		)

		if err != nil {

			http.Error(
				w,
				"Unable to load success page",
				http.StatusInternalServerError,
			)

			return
		}

		tmpl.Execute(
			w,
			verification,
		)

	default:

		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)
	}
}
