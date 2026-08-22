package services

import (
	"errors"
	"sync"
	"time"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/repositories"
)

// VerificationService handles voter verification business rules.
type VerificationService struct {
	verificationRepo repositories.VerificationRepository
	userRepo         repositories.UserRepository

	// Protects approval/rejection operations from concurrent
	// requests inside this application process.
	mu sync.Mutex
}

// NewVerificationService creates a new verification service.
func NewVerificationService(
	verificationRepo repositories.VerificationRepository,
	userRepo repositories.UserRepository,
) *VerificationService {

	return &VerificationService{
		verificationRepo: verificationRepo,
		userRepo:         userRepo,
	}
}

// ============================================================
// Submit Verification
// ============================================================

// SubmitVerification submits a new voter verification request.
func (s *VerificationService) SubmitVerification(
	verification models.Verification,
) error {

	if verification.UserID <= 0 {
		return errors.New("invalid user ID")
	}

	// Make sure the user actually exists.
	_, err := s.userRepo.GetByID(
		verification.UserID,
	)

	if err != nil {
		return errors.New("user not found")
	}

	// Check whether this user already has a verification record.
	existing, err :=
		s.verificationRepo.GetByUserID(
			verification.UserID,
		)

	if err == nil && existing != nil {

		return errors.New(
			"verification already submitted",
		)
	}

	// New verification requests always begin as pending.
	verification.Status = "pending"
	verification.SubmittedAt = time.Now()
	verification.ReviewedAt = time.Time{}
	verification.ReviewedBy = 0

	return s.verificationRepo.Create(
		verification,
	)
}

// ============================================================
// Get User Verification
// ============================================================

// GetUserVerification returns the verification record
// belonging to a user.
func (s *VerificationService) GetUserVerification(
	userID int,
) (*models.Verification, error) {

	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	return s.verificationRepo.GetByUserID(
		userID,
	)
}

// ============================================================
// Get All Verifications
// ============================================================

// GetAllVerifications returns all verification requests.
func (s *VerificationService) GetAllVerifications() []models.Verification {

	return s.verificationRepo.GetAll()
}

// ============================================================
// Approve Verification
// ============================================================

// ApproveVerification approves a pending voter verification.
//
// adminID identifies the administrator who performed the
// approval.
func (s *VerificationService) ApproveVerification(
	verificationID int,
	adminID int,
) error {

	// Prevent concurrent approval/rejection operations
	// inside this application process.
	s.mu.Lock()
	defer s.mu.Unlock()

	if verificationID <= 0 {
		return errors.New(
			"invalid verification ID",
		)
	}

	if adminID <= 0 {
		return errors.New(
			"invalid administrator ID",
		)
	}

	// --------------------------------------------------------
	// Find verification
	// --------------------------------------------------------

	verification, err :=
		s.GetVerificationByID(
			verificationID,
		)

	if err != nil {
		return err
	}

	// --------------------------------------------------------
	// Prevent double review
	// --------------------------------------------------------

	if verification.Status != "pending" {

		return errors.New(
			"verification has already been reviewed",
		)
	}

	// --------------------------------------------------------
	// Load voter
	// --------------------------------------------------------

	user, err :=
		s.userRepo.GetByID(
			verification.UserID,
		)

	if err != nil {
		return errors.New(
			"voter account not found",
		)
	}

	if user == nil {
		return errors.New(
			"voter account not found",
		)
	}

	// --------------------------------------------------------
	// Prevent approving an already verified voter
	// --------------------------------------------------------

	if user.IsVerified {

		return errors.New(
			"user is already verified",
		)
	}

	// --------------------------------------------------------
	// Preserve original records for rollback
	// --------------------------------------------------------

	originalVerification := *verification
	originalUser := *user

	now := time.Now()

	// --------------------------------------------------------
	// Prepare approved verification
	// --------------------------------------------------------

	verification.Status = "approved"
	verification.ReviewedAt = now
	verification.ReviewedBy = adminID

	// --------------------------------------------------------
	// Prepare verified user
	// --------------------------------------------------------

	user.IsVerified = true
	user.AccountActive = true
	user.UpdatedAt = now

	// --------------------------------------------------------
	// Save verification first
	// --------------------------------------------------------

	if err := s.verificationRepo.Update(
		*verification,
	); err != nil {

		return err
	}

	// --------------------------------------------------------
	// Save user
	// --------------------------------------------------------

	if err := s.userRepo.Update(
		*user,
	); err != nil {

		// ----------------------------------------------------
		// ROLLBACK
		//
		// The verification was already changed to approved.
		// Restore it so we do not leave the system in an
		// inconsistent state.
		// ----------------------------------------------------

		rollbackErr :=
			s.verificationRepo.Update(
				originalVerification,
			)

		if rollbackErr != nil {

			return errors.New(
				"approval failed and rollback also failed",
			)
		}

		// Restore user as well in case the repository changed
		// anything before returning its error.
		_ = s.userRepo.Update(
			originalUser,
		)

		return errors.New(
			"unable to verify voter account",
		)
	}

	return nil
}

// ============================================================
// Reject Verification
// ============================================================

// RejectVerification rejects a pending voter verification.
//
// adminID identifies the administrator who performed the
// rejection.
func (s *VerificationService) RejectVerification(
	verificationID int,
	adminID int,
) error {

	// Protect against concurrent review requests.
	s.mu.Lock()
	defer s.mu.Unlock()

	if verificationID <= 0 {
		return errors.New(
			"invalid verification ID",
		)
	}

	if adminID <= 0 {
		return errors.New(
			"invalid administrator ID",
		)
	}

	// --------------------------------------------------------
	// Find verification
	// --------------------------------------------------------

	verification, err :=
		s.GetVerificationByID(
			verificationID,
		)

	if err != nil {
		return err
	}

	// --------------------------------------------------------
	// Prevent double review
	// --------------------------------------------------------

	if verification.Status != "pending" {

		return errors.New(
			"verification has already been reviewed",
		)
	}

	// --------------------------------------------------------
	// Update rejection
	// --------------------------------------------------------

	now := time.Now()

	verification.Status = "rejected"
	verification.ReviewedAt = now
	verification.ReviewedBy = adminID

	return s.verificationRepo.Update(
		*verification,
	)
}

// ============================================================
// Get User By ID
// ============================================================

// GetUserByID returns a user by ID.
func (s *VerificationService) GetUserByID(
	id int,
) (*models.User, error) {

	if id <= 0 {
		return nil, errors.New(
			"invalid user ID",
		)
	}

	return s.userRepo.GetByID(id)
}

// ============================================================
// Get Verification By ID
// ============================================================

// GetVerificationByID returns a verification record by ID.
func (s *VerificationService) GetVerificationByID(
	id int,
) (*models.Verification, error) {

	if id <= 0 {
		return nil, errors.New(
			"invalid verification ID",
		)
	}

	verifications :=
		s.verificationRepo.GetAll()

	for i := range verifications {

		if verifications[i].ID == id {

			// Return a copy rather than relying on the
			// range-variable address.
			verification :=
				verifications[i]

			return &verification, nil
		}
	}

	return nil, errors.New(
		"verification not found",
	)
}
