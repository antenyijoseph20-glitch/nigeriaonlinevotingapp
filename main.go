package main

import (
	"fmt"
	"net/http"

	"nigeriaonlinevoting/handlers"
	"nigeriaonlinevoting/middleware"
	"nigeriaonlinevoting/repositories"
	"nigeriaonlinevoting/services"
)

func main() {

	// ==============================
	// Repository Layer
	// ==============================

	userRepo := repositories.NewJSONRepository(
		"data/users.json",
	)

	verificationRepo := repositories.NewVerificationRepository(
		"data/verifications.json",
	)
	electionRepo := repositories.NewElectionRepository(
		"data/elections.json",
	)

	// ==============================
	// Service Layer
	// ==============================

	authService := services.NewAuthService(
		userRepo,
	)

	verificationService := services.NewVerificationService(
		verificationRepo,
		userRepo,
	)
	electionService := services.NewElectionService(
		electionRepo,
	)

	adminService := services.NewAdminService(
		userRepo,
		verificationRepo,
	)

	// ==============================
	// Handler Layer
	// ==============================

	authHandler := handlers.NewAuthHandler(
		authService,
	)

	dashboardHandler := handlers.NewDashboardHandler(
		authService,
		verificationService,
	)

	profileHandler := handlers.NewProfileHandler(
		authService,
	)

	verificationHandler := handlers.NewVerificationHandler(
		verificationService,
	)
	electionHandler := handlers.NewElectionHandler(
		electionService,
	)

	adminHandler := handlers.NewAdminHandler(
		adminService,
	)

	adminVerificationHandler := handlers.NewAdminVerificationHandler(
		verificationService,
	)

	adminApproveHandler := handlers.NewAdminApproveHandler(
		verificationService,
	)

	adminRejectHandler := handlers.NewAdminRejectHandler(
		verificationService,
	)

	// ==============================
	// Public Routes
	// ==============================

	http.HandleFunc(
		"/",
		handlers.HomeHandler,
	)

	http.HandleFunc(
		"/register",
		authHandler.Register,
	)

	http.HandleFunc(
		"/login",
		authHandler.Login,
	)

	http.HandleFunc(
		"/logout",
		handlers.Logout,
	)

	// ==============================
	// Protected Routes
	// ==============================

	http.HandleFunc(
		"/dashboard",
		middleware.RequireAuth(
			dashboardHandler.Dashboard,
		),
	)

	http.HandleFunc(
		"/profile",
		middleware.RequireAuth(
			profileHandler.Profile,
		),
	)

	http.HandleFunc(
		"/verification",
		middleware.RequireAuth(
			verificationHandler.Verification,
		),
	)

	// ==============================
	// Admin Routes
	// ==============================

	http.HandleFunc(
		"/admin",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				adminHandler.Dashboard,
			),
		),
	)

	http.HandleFunc(
		"/admin/verification",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				adminVerificationHandler.View,
			),
		),
	)

	http.HandleFunc(
		"/admin/approve",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				adminApproveHandler.Approve,
			),
		),
	)

	http.HandleFunc(
		"/admin/reject",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				adminRejectHandler.Reject,
			),
		),
	)
	http.HandleFunc(
		"/admin/elections",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				electionHandler.ElectionDashboard,
			),
		),
	)

	// ==============================
	// Static Files
	// ==============================

	fs := http.FileServer(
		http.Dir("./static"),
	)

	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			fs,
		),
	)

	// ==============================
	// Start Server
	// ==============================

	fmt.Println("=====================================")
	fmt.Println(" 🇳🇬 Nigeria Online Voting System")
	fmt.Println(" Server running on http://localhost:8080")
	fmt.Println("=====================================")

	err := http.ListenAndServe(
		":8080",
		nil,
	)

	if err != nil {
		panic(err)
	}
}
