package handlers

import (
	"net/http"
	"strconv"

	"nigeriaonlinevoting/services"
)

type AdminApproveHandler struct {
	VerificationService *services.VerificationService
}

func NewAdminApproveHandler(
	service *services.VerificationService,
) *AdminApproveHandler {

	return &AdminApproveHandler{
		VerificationService: service,
	}
}

func (h *AdminApproveHandler) Approve(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {

		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	idString := r.FormValue("verificationID")

	id, err := strconv.Atoi(idString)

	if err != nil {

		http.Error(
			w,
			"Invalid verification ID.",
			http.StatusBadRequest,
		)

		return
	}

	err = h.VerificationService.ApproveVerification(id)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	http.Redirect(
		w,
		r,
		"/admin",
		http.StatusSeeOther,
	)
}
