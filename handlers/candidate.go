package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/services"
)

type CandidateHandler struct {
	CandidateService *services.CandidateService
	ElectionService  *services.ElectionService
	PartyService     *services.PartyService
	PositionService  *services.PositionService
}

func NewCandidateHandler(
	candidateService *services.CandidateService,
	electionService *services.ElectionService,
	partyService *services.PartyService,
	positionService *services.PositionService,
) *CandidateHandler {

	return &CandidateHandler{
		CandidateService: candidateService,
		ElectionService:  electionService,
		PartyService:     partyService,
		PositionService:  positionService,
	}
}

type CandidatePageData struct {
	Candidates []models.CandidateView

	Elections []models.Election

	Parties []models.Party

	Positions []models.Position
}

func (h *CandidateHandler) CandidateDashboard(
	w http.ResponseWriter,
	r *http.Request,
) {

	switch r.Method {

	case http.MethodGet:

		data := CandidatePageData{

			Candidates: h.CandidateService.GetAllCandidateViews(),

			Elections: h.ElectionService.GetAllElections(),

			Parties: h.PartyService.GetAllParties(),

			Positions: h.PositionService.GetAllPositions(),
		}

		tmpl, err := template.ParseFiles(
			"templates/candidates.html",
		)

		if err != nil {

			http.Error(
				w,
				"Unable to load candidate page.",
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
				"Unable to render candidate page.",
				http.StatusInternalServerError,
			)

			return
		}

	case http.MethodPost:

		electionID, err := strconv.Atoi(
			r.FormValue("election_id"),
		)

		if err != nil {

			http.Error(
				w,
				"Invalid election.",
				http.StatusBadRequest,
			)

			return
		}

		partyID, err := strconv.Atoi(
			r.FormValue("party_id"),
		)

		if err != nil {

			http.Error(
				w,
				"Invalid political party.",
				http.StatusBadRequest,
			)

			return
		}

		positionID, err := strconv.Atoi(
			r.FormValue("position_id"),
		)

		if err != nil {

			http.Error(
				w,
				"Invalid position.",
				http.StatusBadRequest,
			)

			return
		}

		candidate := models.Candidate{

			ElectionID: electionID,

			PartyID: partyID,

			PositionID: positionID,

			FirstName: r.FormValue("first_name"),

			LastName: r.FormValue("last_name"),

			Gender: r.FormValue("gender"),

			DateOfBirth: r.FormValue("date_of_birth"),

			Email: r.FormValue("email"),

			PhoneNumber: r.FormValue("phone_number"),

			Biography: r.FormValue("biography"),

			Manifesto: r.FormValue("manifesto"),

			Photo: r.FormValue("photo"),
		}

		err = h.CandidateService.CreateCandidate(
			candidate,
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
			"/admin/candidates",
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

func (h *CandidateHandler) Approve(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := strconv.Atoi(
		r.URL.Query().Get("id"),
	)

	if err != nil {

		http.Error(
			w,
			"Invalid candidate ID.",
			http.StatusBadRequest,
		)

		return
	}

	err = h.CandidateService.ApproveCandidate(
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
		"/admin/candidates",
		http.StatusSeeOther,
	)
}

func (h *CandidateHandler) Activate(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := strconv.Atoi(
		r.URL.Query().Get("id"),
	)

	if err != nil {

		http.Error(
			w,
			"Invalid candidate ID.",
			http.StatusBadRequest,
		)

		return
	}

	err = h.CandidateService.ActivateCandidate(
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
		"/admin/candidates",
		http.StatusSeeOther,
	)
}

func (h *CandidateHandler) Deactivate(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := strconv.Atoi(
		r.URL.Query().Get("id"),
	)

	if err != nil {

		http.Error(
			w,
			"Invalid candidate ID.",
			http.StatusBadRequest,
		)

		return
	}

	err = h.CandidateService.DeactivateCandidate(
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
		"/admin/candidates",
		http.StatusSeeOther,
	)
}

func (h *CandidateHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := strconv.Atoi(
		r.URL.Query().Get("id"),
	)

	if err != nil {

		http.Error(
			w,
			"Invalid candidate ID.",
			http.StatusBadRequest,
		)

		return
	}

	err = h.CandidateService.DeleteCandidate(
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
		"/admin/candidates",
		http.StatusSeeOther,
	)
}
