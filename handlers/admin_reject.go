package handlers

import (
	"net/http"
	"strconv"

	"nigeriaonlinevoting/services"
)

type AdminRejectHandler struct {
	VerificationService *services.VerificationService
}

func NewAdminRejectHandler(
	service *services.VerificationService,
) *AdminRejectHandler {

	return &AdminRejectHandler{
		VerificationService: service,
	}
}

func (h *AdminRejectHandler) Reject(
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

	err = h.VerificationService.RejectVerification(id)

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
