package services

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/repositories"
)

type AuthService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo: repo,
	}
}

// Register creates a new user after applying business rules.
func (a *AuthService) Register(user models.User) error {

	// Check if the email already exists.
	existingUser, _ := a.userRepo.GetByEmail(user.Email)
	if existingUser != nil {
		return errors.New("an account with this email already exists")
	}

	// Hash the password.
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.PasswordHash),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)

	// Default values.
	user.Role = "voter"
	user.IsVerified = false
	user.HasVoted = false

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return a.userRepo.Create(user)
}

// Login authenticates a user.
func (a *AuthService) Login(email, password string) (*models.User, error) {

	user, err := a.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

// GetUserByID returns one user.
func (a *AuthService) GetUserByID(id int) (*models.User, error) {
	return a.userRepo.GetByID(id)
}

// GetAllUsers returns all registered users.
func (a *AuthService) GetAllUsers() []models.User {
	return a.userRepo.GetAll()
}

// UpdateUser updates a user's information.
func (a *AuthService) UpdateUser(user models.User) error {
	user.UpdatedAt = time.Now()
	return a.userRepo.Update(user)
}

// DeleteUser removes a user.
func (a *AuthService) DeleteUser(id int) error {
	return a.userRepo.Delete(id)
}
