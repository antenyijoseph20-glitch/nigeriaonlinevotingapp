package handlers

import (
	"html/template"
	"net/http"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/services"
)

type AdminHandler struct {
	AdminService *services.AdminService
}

func NewAdminHandler(
	service *services.AdminService,
) *AdminHandler {

	return &AdminHandler{
		AdminService: service,
	}
}

type AdminDashboardData struct {
	Statistics    models.AdminStatistics
	Verifications []models.Verification
}

func (h *AdminHandler) Dashboard(
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

	data := AdminDashboardData{

		Statistics: h.AdminService.GetDashboardStatistics(),

		Verifications: h.AdminService.GetPendingVerifications(),
	}

	tmpl, err := template.ParseFiles(
		"templates/admin_dashboard.html",
	)

	if err != nil {

		http.Error(
			w,
			"Unable to load admin dashboard.",
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
			"Unable to render admin dashboard.",
			http.StatusInternalServerError,
		)

		return
	}
}
