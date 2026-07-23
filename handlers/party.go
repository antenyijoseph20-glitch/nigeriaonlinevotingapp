package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/services"
)

type PartyHandler struct {
	PartyService *services.PartyService
}

func NewPartyHandler(
	service *services.PartyService,
) *PartyHandler {

	return &PartyHandler{
		PartyService: service,
	}
}

type PartyPageData struct {
	Parties []models.Party
}

// ===================================
// Party Dashboard
// ===================================

func (h *PartyHandler) PartyDashboard(
	w http.ResponseWriter,
	r *http.Request,
) {

	switch r.Method {

	case http.MethodGet:

		data := PartyPageData{

			Parties: h.PartyService.GetAllParties(),
		}

		tmpl, err := template.ParseFiles(
			"templates/parties.html",
		)

		if err != nil {

			http.Error(
				w,
				"Unable to load parties page.",
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
				"Unable to render parties page.",
				http.StatusInternalServerError,
			)

			return
		}

	case http.MethodPost:

		party := models.Party{

			Name:         r.FormValue("name"),
			Abbreviation: r.FormValue("abbreviation"),
			Slogan:       r.FormValue("slogan"),
			Chairman:     r.FormValue("chairman"),
			Headquarters: r.FormValue("headquarters"),
			Description:  r.FormValue("description"),
			Logo:         r.FormValue("logo"),
		}

		err := h.PartyService.CreateParty(
			party,
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
			"/admin/parties",
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

// ===================================
// Activate Party
// ===================================

func (h *PartyHandler) Activate(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := strconv.Atoi(
		r.URL.Query().Get("id"),
	)

	if err != nil {

		http.Error(
			w,
			"Invalid Party ID",
			http.StatusBadRequest,
		)

		return
	}

	err = h.PartyService.ActivateParty(
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
		"/admin/parties",
		http.StatusSeeOther,
	)
}

// ===================================
// Deactivate Party
// ===================================

func (h *PartyHandler) Deactivate(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := strconv.Atoi(
		r.URL.Query().Get("id"),
	)

	if err != nil {

		http.Error(
			w,
			"Invalid Party ID",
			http.StatusBadRequest,
		)

		return
	}

	err = h.PartyService.DeactivateParty(
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
		"/admin/parties",
		http.StatusSeeOther,
	)
}

// ===================================
// Delete Party
// ===================================

func (h *PartyHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := strconv.Atoi(
		r.URL.Query().Get("id"),
	)

	if err != nil {

		http.Error(
			w,
			"Invalid Party ID",
			http.StatusBadRequest,
		)

		return
	}

	err = h.PartyService.DeleteParty(
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
		"/admin/parties",
		http.StatusSeeOther,
	)
}
