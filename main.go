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

	// ==============================
	// Handler Layer
	// ==============================

	authHandler := handlers.NewAuthHandler(
		authService,
	)

	dashboardHandler := handlers.NewDashboardHandler(
		authService,
	)

	profileHandler := handlers.NewProfileHandler(
		authService,
	)

	verificationHandler := handlers.NewVerificationHandler(
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
