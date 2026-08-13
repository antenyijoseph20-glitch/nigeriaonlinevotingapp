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

	// =====================================
	// Repository Layer
	// =====================================

	userRepo := repositories.NewJSONRepository(
		"data/users.json",
	)

	verificationRepo := repositories.NewVerificationRepository(
		"data/verifications.json",
	)

	electionRepo := repositories.NewElectionRepository(
		"data/elections.json",
	)

	partyRepo := repositories.NewPartyRepository(
		"data/parties.json",
	)

	positionRepo := repositories.NewPositionRepository(
		"data/positions.json",
	)

	candidateRepo := repositories.NewCandidateRepository(
		"data/candidates.json",
	)

	// =====================================
	// Service Layer
	// =====================================

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

	partyService := services.NewPartyService(
		partyRepo,
	)

	positionService := services.NewPositionService(
		positionRepo,
	)

	candidateService := services.NewCandidateService(
		candidateRepo,
		electionRepo,
		partyRepo,
		positionRepo,
	)

	adminService := services.NewAdminService(
		userRepo,
		verificationRepo,
	)

	// =====================================
	// Handler Layer
	// =====================================

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

	electionHandler := handlers.NewElectionHandler(
		electionService,
	)

	partyHandler := handlers.NewPartyHandler(
		partyService,
	)

	positionHandler := handlers.NewPositionHandler(
		positionService,
	)

	candidateHandler := handlers.NewCandidateHandler(
		candidateService,
		electionService,
		partyService,
		positionService,
	)

	// =====================================
	// Public Routes
	// =====================================

	http.HandleFunc("/", handlers.HomeHandler)

	http.HandleFunc("/register", authHandler.Register)

	http.HandleFunc("/login", authHandler.Login)

	http.HandleFunc("/logout", handlers.Logout)

	// =====================================
	// User Routes
	// =====================================

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

	// =====================================
	// Admin Dashboard
	// =====================================

	http.HandleFunc(
		"/admin",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				adminHandler.Dashboard,
			),
		),
	)

	// =====================================
	// Verification
	// =====================================

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

	// =====================================
	// Elections
	// =====================================

	http.HandleFunc(
		"/admin/elections",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				electionHandler.ElectionDashboard,
			),
		),
	)

	http.HandleFunc(
		"/admin/election/open",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				electionHandler.OpenElection,
			),
		),
	)

	http.HandleFunc(
		"/admin/election/close",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				electionHandler.CloseElection,
			),
		),
	)

	http.HandleFunc(
		"/admin/election/delete",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				electionHandler.DeleteElection,
			),
		),
	)

	// =====================================
	// Political Parties
	// =====================================

	http.HandleFunc(
		"/admin/parties",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				partyHandler.PartyDashboard,
			),
		),
	)

	http.HandleFunc(
		"/admin/party/activate",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				partyHandler.Activate,
			),
		),
	)

	http.HandleFunc(
		"/admin/party/deactivate",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				partyHandler.Deactivate,
			),
		),
	)

	http.HandleFunc(
		"/admin/party/delete",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				partyHandler.Delete,
			),
		),
	)

	// =====================================
	// Positions
	// =====================================

	http.HandleFunc(
		"/admin/positions",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				positionHandler.PositionDashboard,
			),
		),
	)

	// =====================================
	// Candidates
	// =====================================

	http.HandleFunc(
		"/admin/candidates",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				candidateHandler.CandidateDashboard,
			),
		),
	)

	http.HandleFunc(
		"/admin/candidate/approve",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				candidateHandler.Approve,
			),
		),
	)

	http.HandleFunc(
		"/admin/candidate/activate",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				candidateHandler.Activate,
			),
		),
	)

	http.HandleFunc(
		"/admin/candidate/deactivate",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				candidateHandler.Deactivate,
			),
		),
	)

	http.HandleFunc(
		"/admin/candidate/delete",
		middleware.RequireAuth(
			middleware.RequireAdmin(
				candidateHandler.Delete,
			),
		),
	)

	// =====================================
	// Static Files
	// =====================================

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

	// =====================================
	// Start Server
	// =====================================

	fmt.Println("=====================================")
	fmt.Println(" 🇳🇬 Nigeria Online Voting System")
	fmt.Println(" Server running on http://localhost:8080")
	fmt.Println("=====================================")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
