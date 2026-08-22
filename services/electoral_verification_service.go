package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/providers"
	"nigeriaonlinevoting/repositories"
)

// ElectoralVerificationService handles verification of a voter
// against an authorized electoral verification provider.
//
// This service deliberately keeps the external provider behind
// the ElectoralVerificationProvider interface.
//
// The application must never assume that a voter is registered
// merely because a provider is unavailable.
type ElectoralVerificationService struct {
	ElectoralVerificationRepo repositories.ElectoralVerificationRepository
	UserRepo                   repositories.UserRepository
	Provider                   providers.ElectoralVerificationProvider
}

// NewElectoralVerificationService creates a new electoral
// verification service.
func NewElectoralVerificationService(
	electoralVerificationRepo repositories.ElectoralVerificationRepository,
	userRepo repositories.UserRepository,
	provider providers.ElectoralVerificationProvider,
) *ElectoralVerificationService {

	return &ElectoralVerificationService{
		ElectoralVerificationRepo: electoralVerificationRepo,
		UserRepo:                   userRepo,
		Provider:                   provider,
	}
}

// VerifyVoter performs an electoral verification request.
//
// The request is sent to the configured provider. The service
// stores the provider result but never invents a successful
// verification when the provider is unavailable.
func (s *ElectoralVerificationService) VerifyVoter(
	ctx context.Context,
	userID int,
) (*models.ElectoralVerification, error) {

	// -------------------------------------
	// Validate user ID
	// -------------------------------------

	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	// -------------------------------------
	// Load user
	// -------------------------------------

	user, err := s.UserRepo.GetByID(userID)

	if err != nil {
		return nil, errors.New("user not found")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// -------------------------------------
	// Check existing electoral verification
	// -------------------------------------

	existing, err :=
		s.ElectoralVerificationRepo.GetByUserID(userID)

	if err == nil && existing != nil {

		switch existing.Status {

		case "verified":
			return existing, errors.New(
				"voter has already been electorally verified",
			)

		case "pending":
			return existing, errors.New(
				"electoral verification is already pending",
			)
		}
	}

	// -------------------------------------
	// Validate required identity information
	// -------------------------------------

	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.NIN = strings.TrimSpace(user.NIN)
	user.VIN = strings.TrimSpace(user.VIN)

	if user.FirstName == "" {
		return nil, errors.New("first name is required")
	}

	if user.LastName == "" {
		return nil, errors.New("last name is required")
	}

	if user.VIN == "" {
		return nil, errors.New("VIN is required")
	}

	if user.NIN == "" {
		return nil, errors.New("NIN is required")
	}

	// -------------------------------------
	// Build provider request
	// -------------------------------------

	request := providers.VoterVerificationRequest{
		VIN:         user.VIN,
		NIN:         user.NIN,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DateOfBirth: user.DateOfBirth.Format("2006-01-02"),
		State:       user.State,
		LGA:         user.LGA,
	}

	// -------------------------------------
	// Ask the configured provider
	// -------------------------------------

	result, providerErr :=
		s.Provider.VerifyVoter(
			ctx,
			request,
		)

	now := time.Now()

	verification := models.ElectoralVerification{
		UserID:           userID,
		ProviderName:     strings.TrimSpace(result.ProviderName),
		ReferenceID:      strings.TrimSpace(result.ReferenceID),
		RegisteredVoter:  result.RegisteredVoter,
		VINMatched:       result.VINMatched,
		IdentityMatched:  result.IdentityMatched,
		Message:          strings.TrimSpace(result.Message),
		RequestedAt:      now,
		VerifiedAt:      now,
	}

	// -------------------------------------
	// Provider unavailable/error
	// -------------------------------------

	if providerErr != nil {

		verification.Status = "unavailable"

		// An unavailable provider must never produce
		// a successful electoral verification.
		verification.RegisteredVoter = false
		verification.VINMatched = false
		verification.IdentityMatched = false

		if verification.ProviderName == "" {
			verification.ProviderName = "unavailable"
		}

		if verification.Message == "" {
			verification.Message =
				"Electoral verification provider is unavailable."
		}

		if err := s.ElectoralVerificationRepo.Create(
			verification,
		); err != nil {
			return nil, err
		}

		return &verification, providerErr
	}

	// -------------------------------------
	// Successful provider response
	// -------------------------------------

	// A voter is only considered electorally verified
	// when the provider explicitly confirms verification.
	if result.Verified &&
		result.RegisteredVoter &&
		result.VINMatched &&
		result.IdentityMatched {

		verification.Status = "verified"

	} else {

		verification.Status = "failed"

		// Defensive rule:
		// if the provider did not confirm verification,
		// the application must not treat the voter as verified.
		verification.RegisteredVoter =
			result.RegisteredVoter

		verification.VINMatched =
			result.VINMatched

		verification.IdentityMatched =
			result.IdentityMatched
	}

	// -------------------------------------
	// Store result
	// -------------------------------------

	if err := s.ElectoralVerificationRepo.Create(
		verification,
	); err != nil {
		return nil, err
	}

	return &verification, nil
}

// GetByUserID returns the electoral verification record
// belonging to a user.
func (s *ElectoralVerificationService) GetByUserID(
	userID int,
) (*models.ElectoralVerification, error) {

	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	return s.ElectoralVerificationRepo.GetByUserID(
		userID,
	)
}

// GetByID returns an electoral verification record by ID.
func (s *ElectoralVerificationService) GetByID(
	id int,
) (*models.ElectoralVerification, error) {

	if id <= 0 {
		return nil, errors.New(
			"invalid electoral verification ID",
		)
	}

	return s.ElectoralVerificationRepo.GetByID(id)
}

// GetAll returns all electoral verification records.
func (s *ElectoralVerificationService) GetAll() []models.ElectoralVerification {

	return s.ElectoralVerificationRepo.GetAll()
}