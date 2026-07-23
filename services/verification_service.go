package services

import (
	"errors"
	"time"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/repositories"
)

type VerificationService struct {
	verificationRepo repositories.VerificationRepository

	userRepo repositories.UserRepository
}

// Constructor

func NewVerificationService(
	verificationRepo repositories.VerificationRepository,
	userRepo repositories.UserRepository,
) *VerificationService {

	return &VerificationService{

		verificationRepo: verificationRepo,

		userRepo: userRepo,
	}
}

// Submit a voter verification request

func (s *VerificationService) SubmitVerification(
	verification models.Verification,
) error {

	// Check if user already submitted verification

	existing, err := s.verificationRepo.GetByUserID(
		verification.UserID,
	)

	if err == nil && existing != nil {

		return errors.New(
			"verification already submitted",
		)
	}

	// Default verification status

	verification.Status = "pending"

	verification.SubmittedAt = time.Now()

	return s.verificationRepo.Create(
		verification,
	)

}

// Get verification record for a user

func (s *VerificationService) GetUserVerification(
	userID int,
) (*models.Verification, error) {

	return s.verificationRepo.GetByUserID(
		userID,
	)

}

// Get all verification requests
// Mainly for admin dashboard

func (s *VerificationService) GetAllVerifications() []models.Verification {

	return s.verificationRepo.GetAll()

}

// Approve voter verification

func (s *VerificationService) ApproveVerification(
	verificationID int,
) error {

	// Find verification record

	verifications :=
		s.verificationRepo.GetAll()

	for _, verification := range verifications {

		if verification.ID == verificationID {

			// Update verification status

			verification.Status = "approved"

			verification.ReviewedAt = time.Now()

			err := s.verificationRepo.Update(
				verification,
			)

			if err != nil {

				return err

			}

			// Update user account

			user, err :=
				s.userRepo.GetByID(
					verification.UserID,
				)

			if err != nil {

				return err

			}

			user.IsVerified = true

			user.UpdatedAt = time.Now()

			return s.userRepo.Update(
				*user,
			)

		}
	}

	return errors.New(
		"verification record not found",
	)

}

// Reject verification

func (s *VerificationService) RejectVerification(
	verificationID int,
) error {

	verifications :=
		s.verificationRepo.GetAll()

	for _, verification := range verifications {

		if verification.ID == verificationID {

			verification.Status = "rejected"

			verification.ReviewedAt = time.Now()

			return s.verificationRepo.Update(
				verification,
			)

		}
	}

	return errors.New(
		"verification record not found",
	)

}

// GetUserByID returns a user by ID.
func (s *VerificationService) GetUserByID(
	id int,
) (*models.User, error) {

	return s.userRepo.GetByID(id)
}

// GetVerificationByID returns a verification record by its ID.
func (s *VerificationService) GetVerificationByID(
	id int,
) (*models.Verification, error) {

	verifications := s.verificationRepo.GetAll()

	for _, verification := range verifications {

		if verification.ID == id {
			return &verification, nil
		}
	}

	return nil, errors.New("verification not found")
}
