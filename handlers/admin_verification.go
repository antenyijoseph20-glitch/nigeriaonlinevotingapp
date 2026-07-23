package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"nigeriaonlinevoting/services"
)

type AdminVerificationHandler struct {
	VerificationService *services.VerificationService
}

func NewAdminVerificationHandler(
	service *services.VerificationService,
) *AdminVerificationHandler {

	return &AdminVerificationHandler{
		VerificationService: service,
	}
}

func (h *AdminVerificationHandler) View(
	w http.ResponseWriter,
	r *http.Request,
) {

	idString := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idString)
	if err != nil {

		http.Error(
			w,
			"Invalid verification ID.",
			http.StatusBadRequest,
		)

		return
	}

	verification, err :=
		h.VerificationService.GetVerificationByID(id)

	if err != nil {

		http.Error(
			w,
			"Verification not found.",
			http.StatusNotFound,
		)

		return
	}

	tmpl, err := template.ParseFiles(
		"templates/admin_verification.html",
	)

	if err != nil {

		http.Error(
			w,
			"Unable to load page.",
			http.StatusInternalServerError,
		)

		return
	}

	tmpl.Execute(
		w,
		verification,
	)
}
