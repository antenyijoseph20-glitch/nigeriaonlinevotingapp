package handlers

import (
	"bytes"
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

	// =====================================
	// HTTP METHOD CHECK
	// =====================================

	if r.Method != http.MethodGet {

		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	// =====================================
	// LOAD DASHBOARD DATA
	// =====================================

	data := AdminDashboardData{
		Statistics: h.AdminService.GetDashboardStatistics(),

		Verifications: h.AdminService.GetPendingVerifications(),
	}

	// =====================================
	// LOAD TEMPLATE
	// =====================================

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

	// =====================================
	// RENDER TEMPLATE INTO BUFFER
	// =====================================
	//
	// We render into a buffer first.
	//
	// This prevents a partially-written
	// HTTP response if template execution
	// fails.

	var buffer bytes.Buffer

	err = tmpl.Execute(
		&buffer,
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

	// =====================================
	// SEND SUCCESSFUL RESPONSE
	// =====================================

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	w.WriteHeader(
		http.StatusOK,
	)

	_, err = w.Write(
		buffer.Bytes(),
	)

	if err != nil {
		return
	}
}
