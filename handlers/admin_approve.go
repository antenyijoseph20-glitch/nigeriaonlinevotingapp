package handlers

import (
	"net/http"
	"strconv"

	"nigeriaonlinevoting/services"
	"nigeriaonlinevoting/sessions"
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

	// Get the authenticated administrator's ID
	// from the existing session.
	adminID, err := sessions.GetSessionUserID(r)

	if err != nil {
		http.Error(
			w,
			"Unauthorized.",
			http.StatusUnauthorized,
		)
		return
	}

	// Get verification ID.
	idString := r.FormValue("verificationID")

	id, err := strconv.Atoi(idString)

	if err != nil || id <= 0 {
		http.Error(
			w,
			"Invalid verification ID.",
			http.StatusBadRequest,
		)
		return
	}

	// Approve verification and record
	// which administrator performed the action.
	err = h.VerificationService.ApproveVerification(
		id,
		adminID,
	)

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
