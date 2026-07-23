package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/services"
)

type ElectionHandler struct {
	ElectionService *services.ElectionService
}

func NewElectionHandler(
	service *services.ElectionService,
) *ElectionHandler {

	return &ElectionHandler{
		ElectionService: service,
	}
}

type ElectionPageData struct {
	Elections []models.Election
}

// ElectionDashboard displays all elections and handles creating a new one.
func (h *ElectionHandler) ElectionDashboard(
	w http.ResponseWriter,
	r *http.Request,
) {

	switch r.Method {

	case http.MethodGet:

		data := ElectionPageData{
			Elections: h.ElectionService.GetAllElections(),
		}

		tmpl, err := template.ParseFiles(
			"templates/elections.html",
		)

		if err != nil {

			http.Error(
				w,
				"Unable to load elections page.",
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
				"Unable to render elections page.",
				http.StatusInternalServerError,
			)

			return
		}

	case http.MethodPost:

		title := r.FormValue("title")

		description := r.FormValue("description")

		startDate := r.FormValue("startDate")

		endDate := r.FormValue("endDate")

		start, err := time.Parse(
			"2006-01-02",
			startDate,
		)

		if err != nil {

			http.Error(
				w,
				"Invalid start date.",
				http.StatusBadRequest,
			)

			return
		}

		end, err := time.Parse(
			"2006-01-02",
			endDate,
		)

		if err != nil {

			http.Error(
				w,
				"Invalid end date.",
				http.StatusBadRequest,
			)

			return
		}

		election := models.Election{

			Title: title,

			Description: description,

			StartDate: start,

			EndDate: end,
		}

		err = h.ElectionService.CreateElection(
			election,
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
			"/admin/elections",
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
func (h *ElectionHandler) OpenElection(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodGet {

		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	id, err := strconv.Atoi(
		r.URL.Query().Get("id"),
	)

	if err != nil {

		http.Error(
			w,
			"Invalid election ID.",
			http.StatusBadRequest,
		)

		return
	}

	err = h.ElectionService.OpenElection(id)

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
		"/admin/elections",
		http.StatusSeeOther,
	)
}
func (h *ElectionHandler) CloseElection(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodGet {

		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	id, err := strconv.Atoi(
		r.URL.Query().Get("id"),
	)

	if err != nil {

		http.Error(
			w,
			"Invalid election ID.",
			http.StatusBadRequest,
		)

		return
	}

	err = h.ElectionService.CloseElection(id)

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
		"/admin/elections",
		http.StatusSeeOther,
	)
}
func (h *ElectionHandler) DeleteElection(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodGet {

		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	id, err := strconv.Atoi(
		r.URL.Query().Get("id"),
	)

	if err != nil {

		http.Error(
			w,
			"Invalid election ID.",
			http.StatusBadRequest,
		)

		return
	}

	err = h.ElectionService.DeleteElection(id)

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
		"/admin/elections",
		http.StatusSeeOther,
	)
}
