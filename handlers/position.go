package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/services"
)

type PositionHandler struct {
	PositionService *services.PositionService
}

func NewPositionHandler(
	service *services.PositionService,
) *PositionHandler {

	return &PositionHandler{
		PositionService: service,
	}
}

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

	switch r.Method {

	case http.MethodGet:

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

		err = tmpl.Execute(
			w,
			data,
		)

		if err != nil {

			http.Error(
				w,
				"Unable to render positions page.",
				http.StatusInternalServerError,
			)

			return
		}

	case http.MethodPost:

		seats, err := strconv.Atoi(
			r.FormValue("seats"),
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

			Name: r.FormValue("name"),

			Description: r.FormValue("description"),

			Level: r.FormValue("level"),

			Seats: seats,
		}

		err = h.PositionService.CreatePosition(
			position,
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
			"/admin/positions",
			http.StatusSeeOther,
		)

	default:

		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

// =====================================
// Activate Position
// =====================================

func (h *PositionHandler) Activate(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := strconv.Atoi(
		r.URL.Query().Get("id"),
	)

	if err != nil {

		http.Error(
			w,
			"Invalid Position ID",
			http.StatusBadRequest,
		)

		return
	}

	err = h.PositionService.ActivatePosition(
		id,
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

	id, err := strconv.Atoi(
		r.URL.Query().Get("id"),
	)

	if err != nil {

		http.Error(
			w,
			"Invalid Position ID",
			http.StatusBadRequest,
		)

		return
	}

	err = h.PositionService.DeactivatePosition(
		id,
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

	id, err := strconv.Atoi(
		r.URL.Query().Get("id"),
	)

	if err != nil {

		http.Error(
			w,
			"Invalid Position ID",
			http.StatusBadRequest,
		)

		return
	}

	err = h.PositionService.DeletePosition(
		id,
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
		"/admin/positions",
		http.StatusSeeOther,
	)
}
