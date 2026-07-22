package handlers

import (
	"html/template"
	"net/http"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/services"
	"nigeriaonlinevoting/sessions"
	"nigeriaonlinevoting/validators"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: service,
	}
}

// Register handles user registration.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:

		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(w, "Unable to load register page", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, nil)

	case http.MethodPost:

		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")
		email := r.FormValue("email")
		password := r.FormValue("password")
		phoneNumber := r.FormValue("phoneNumber")
		nin := r.FormValue("nin")

		// Validate required fields
		if firstName == "" {
			http.Error(w, "First name is required.", http.StatusBadRequest)
			return
		}

		if lastName == "" {
			http.Error(w, "Last name is required.", http.StatusBadRequest)
			return
		}

		if email == "" {
			http.Error(w, "Email is required.", http.StatusBadRequest)
			return
		}

		if phoneNumber == "" {
			http.Error(w, "Phone number is required.", http.StatusBadRequest)
			return
		}

		if nin == "" {
			http.Error(w, "NIN is required.", http.StatusBadRequest)
			return
		}

		validEmail, msg := validators.ValidateEmail(email)
		if !validEmail {
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		validPassword, msg := validators.ValidatePassword(password)
		if !validPassword {
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		user := models.User{
			FirstName:    firstName,
			LastName:     lastName,
			Email:        email,
			PhoneNumber:  phoneNumber,
			NIN:          nin,
			PasswordHash: password,
		}

		err := h.AuthService.Register(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := struct {
			FirstName   string
			LastName    string
			Email       string
			PhoneNumber string
			NIN         string
		}{
			FirstName:   firstName,
			LastName:    lastName,
			Email:       email,
			PhoneNumber: phoneNumber,
			NIN:         nin,
		}

		tmpl, err := template.ParseFiles("templates/success.html")
		if err != nil {
			http.Error(w, "Unable to load success page", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Login handles user login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:

		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Unable to load login page", http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "Unable to render login page", http.StatusInternalServerError)
			return
		}

	case http.MethodPost:

		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Error(w, "Email and password are required.", http.StatusBadRequest)
			return
		}

		// Authenticate user
		user, err := h.AuthService.Login(email, password)
		if err != nil {
			http.Error(w, "Invalid email or password.", http.StatusUnauthorized)
			return
		}

		// Create a session cookie
		sessions.CreateSession(w, user.ID)

		// Redirect to dashboard
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
