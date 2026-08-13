package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/services"
)

type PositionHandler struct {
	PositionService *services.PositionService
}

// =====================================
// Constructor
// =====================================

func NewPositionHandler(
	service *services.PositionService,
) *PositionHandler {

	return &PositionHandler{
		PositionService: service,
	}
}

// =====================================
// Page Data
// =====================================

type PositionPageData struct {
	Positions []models.Position
}

// =====================================
// Position Dashboard
// =====================================

func (h *PositionHandler) PositionDashboard(
	w http.ResponseWriter,
	r *http.Request,
) {

	if h.PositionService == nil {
		http.Error(
			w,
			"Position service is unavailable.",
			http.StatusInternalServerError,
		)

		return
	}

	switch r.Method {

	case http.MethodGet:

		h.showPositionDashboard(
			w,
			r,
		)

	case http.MethodPost:

		h.createPosition(
			w,
			r,
		)

	default:

		w.Header().Set(
			"Allow",
			"GET, POST",
		)

		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

// =====================================
// Show Position Dashboard
// =====================================

func (h *PositionHandler) showPositionDashboard(
	w http.ResponseWriter,
	r *http.Request,
) {

	data := PositionPageData{
		Positions: h.PositionService.GetAllPositions(),
	}

	tmpl, err := template.ParseFiles(
		"templates/positions.html",
	)

	if err != nil {

		http.Error(
			w,
			"Unable to load positions page.",
			http.StatusInternalServerError,
		)

		return
	}

	if err := tmpl.Execute(
		w,
		data,
	); err != nil {

		http.Error(
			w,
			"Unable to render positions page.",
			http.StatusInternalServerError,
		)

		return
	}
}

// =====================================
// Create Position
// =====================================

func (h *PositionHandler) createPosition(
	w http.ResponseWriter,
	r *http.Request,
) {

	if err := r.ParseForm(); err != nil {

		http.Error(
			w,
			"Invalid form submission.",
			http.StatusBadRequest,
		)

		return
	}

	seatsText := strings.TrimSpace(
		r.FormValue("seats"),
	)

	if seatsText == "" {

		http.Error(
			w,
			"Number of seats is required.",
			http.StatusBadRequest,
		)

		return
	}

	seats, err := strconv.Atoi(
		seatsText,
	)

	if err != nil {

		http.Error(
			w,
			"Invalid number of seats.",
			http.StatusBadRequest,
		)

		return
	}

	position := models.Position{
		Name: strings.TrimSpace(
			r.FormValue("name"),
		),

		Description: strings.TrimSpace(
			r.FormValue("description"),
		),

		Level: strings.TrimSpace(
			r.FormValue("level"),
		),

		Seats: seats,
	}

	if err := h.PositionService.CreatePosition(
		position,
	); err != nil {

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
		"/admin/positions",
		http.StatusSeeOther,
	)
}

// =====================================
// Activate Position
// =====================================

func (h *PositionHandler) Activate(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {

		w.Header().Set(
			"Allow",
			http.MethodPost,
		)

		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	id, err := getPositionID(
		r,
	)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	if err := h.PositionService.ActivatePosition(
		id,
	); err != nil {

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
		"/admin/positions",
		http.StatusSeeOther,
	)
}

// =====================================
// Deactivate Position
// =====================================

func (h *PositionHandler) Deactivate(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {

		w.Header().Set(
			"Allow",
			http.MethodPost,
		)

		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	id, err := getPositionID(
		r,
	)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	if err := h.PositionService.DeactivatePosition(
		id,
	); err != nil {

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
		"/admin/positions",
		http.StatusSeeOther,
	)
}

// =====================================
// Delete Position
// =====================================

func (h *PositionHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {

		w.Header().Set(
			"Allow",
			http.MethodPost,
		)

		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	id, err := getPositionID(
		r,
	)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	if err := h.PositionService.DeletePosition(
		id,
	); err != nil {

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
		"/admin/positions",
		http.StatusSeeOther,
	)
}

// =====================================
// Get Position ID
// =====================================

func getPositionID(
	r *http.Request,
) (int, error) {

	idText := strings.TrimSpace(
		r.URL.Query().Get("id"),
	)

	if idText == "" {
		return 0, strconv.ErrSyntax
	}

	id, err := strconv.Atoi(
		idText,
	)

	if err != nil {
		return 0, err
	}

	if id <= 0 {
		return 0, strconv.ErrSyntax
	}

	return id, nil
}
