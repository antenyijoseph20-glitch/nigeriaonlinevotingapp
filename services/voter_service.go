package services

import (
	"errors"
	"strings"
	"time"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/repositories"
)

// VoterService contains the business rules
// for managing voters.
type VoterService struct {
	VoterRepo repositories.VoterRepository
	UserRepo  repositories.UserRepository
}

// =====================================
// Constructor
// =====================================

func NewVoterService(
	voterRepo repositories.VoterRepository,
	userRepo repositories.UserRepository,
) *VoterService {

	return &VoterService{
		VoterRepo: voterRepo,
		UserRepo:  userRepo,
	}
}

// =====================================
// Register Voter
// =====================================

func (s *VoterService) RegisterVoter(
	voter models.Voter,
) error {

	// -------------------------------------
	// Basic service validation
	// -------------------------------------

	if voter.UserID <= 0 {
		return errors.New(
			"invalid user ID",
		)
	}

	voter.NIN = strings.TrimSpace(voter.NIN)
	voter.VIN = strings.TrimSpace(voter.VIN)
	voter.PVCNumber = strings.TrimSpace(voter.PVCNumber)

	if voter.NIN == "" {
		return errors.New(
			"NIN is required",
		)
	}

	if voter.VIN == "" {
		return errors.New(
			"VIN is required",
		)
	}

	if voter.PVCNumber == "" {
		return errors.New(
			"PVC number is required",
		)
	}

	// -------------------------------------
	// Verify linked user exists
	// -------------------------------------

	_, err := s.UserRepo.GetByID(
		voter.UserID,
	)

	if err != nil {
		return errors.New(
			"linked user not found",
		)
	}

	// -------------------------------------
	// Prevent duplicate voter records
	// -------------------------------------

	_, err = s.VoterRepo.GetByUserID(
		voter.UserID,
	)

	if err == nil {
		return errors.New(
			"voter already exists for this user",
		)
	}

	// -------------------------------------
	// Prevent duplicate NIN
	// -------------------------------------

	_, err = s.VoterRepo.GetByNIN(
		voter.NIN,
	)

	if err == nil {
		return errors.New(
			"NIN already belongs to another voter",
		)
	}

	// -------------------------------------
	// Prevent duplicate VIN
	// -------------------------------------

	_, err = s.VoterRepo.GetByVIN(
		voter.VIN,
	)

	if err == nil {
		return errors.New(
			"VIN already belongs to another voter",
		)
	}

	// -------------------------------------
	// Prevent duplicate PVC
	// -------------------------------------

	_, err = s.VoterRepo.GetByPVCNumber(
		voter.PVCNumber,
	)

	if err == nil {
		return errors.New(
			"PVC number already belongs to another voter",
		)
	}

	// -------------------------------------
	// Initial voter state
	// -------------------------------------

	voter.IsRegistered = true
	voter.IsVerified = false
	voter.IsEligible = false
	voter.IsSuspended = false
	voter.HasVoted = false

	voter.CreatedAt = time.Now()
	voter.UpdatedAt = voter.CreatedAt

	// -------------------------------------
	// Store voter
	// -------------------------------------

	return s.VoterRepo.Create(
		voter,
	)
}

// =====================================
// Get Voter
// =====================================

func (s *VoterService) GetVoter(
	id int,
) (*models.Voter, error) {

	if id <= 0 {
		return nil, errors.New(
			"invalid voter ID",
		)
	}

	voter, err := s.VoterRepo.GetByID(
		id,
	)

	if err != nil {
		return nil, err
	}

	return voter, nil
}

// =====================================
// Get Voter By User
// =====================================

func (s *VoterService) GetVoterByUserID(
	userID int,
) (*models.Voter, error) {

	if userID <= 0 {
		return nil, errors.New(
			"invalid user ID",
		)
	}

	return s.VoterRepo.GetByUserID(
		userID,
	)
}

// =====================================
// Verify Voter
// =====================================

func (s *VoterService) VerifyVoter(
	id int,
) error {

	if id <= 0 {
		return errors.New(
			"invalid voter ID",
		)
	}

	voter, err := s.VoterRepo.GetByID(
		id,
	)

	if err != nil {
		return err
	}

	if !voter.IsRegistered {
		return errors.New(
			"voter is not registered",
		)
	}

	if voter.IsSuspended {
		return errors.New(
			"suspended voter cannot be verified",
		)
	}

	if voter.AccountLocked {
		return errors.New(
			"locked voter cannot be verified",
		)
	}

	voter.IsVerified = true

	// Eligibility requires successful
	// registration and verification.
	voter.IsEligible = voter.IsRegistered &&
		voter.IsVerified &&
		!voter.IsSuspended

	voter.LastVerification = time.Now()
	voter.UpdatedAt = time.Now()

	return s.VoterRepo.Update(
		*voter,
	)
}

// =====================================
// Suspend Voter
// =====================================

func (s *VoterService) SuspendVoter(
	id int,
) error {

	if id <= 0 {
		return errors.New(
			"invalid voter ID",
		)
	}

	voter, err := s.VoterRepo.GetByID(
		id,
	)

	if err != nil {
		return err
	}

	if voter.IsSuspended {
		return errors.New(
			"voter is already suspended",
		)
	}

	voter.IsSuspended = true
	voter.IsEligible = false
	voter.UpdatedAt = time.Now()

	return s.VoterRepo.Update(
		*voter,
	)
}

// =====================================
// Reinstate Voter
// =====================================

func (s *VoterService) ReinstateVoter(
	id int,
) error {

	if id <= 0 {
		return errors.New(
			"invalid voter ID",
		)
	}

	voter, err := s.VoterRepo.GetByID(
		id,
	)

	if err != nil {
		return err
	}

	if !voter.IsSuspended {
		return errors.New(
			"voter is not suspended",
		)
	}

	voter.IsSuspended = false

	voter.IsEligible =
		voter.IsRegistered &&
			voter.IsVerified

	voter.UpdatedAt = time.Now()

	return s.VoterRepo.Update(
		*voter,
	)
}

// =====================================
// Enroll Face
// =====================================

func (s *VoterService) EnrollFace(
	id int,
) error {

	if id <= 0 {
		return errors.New(
			"invalid voter ID",
		)
	}

	voter, err := s.VoterRepo.GetByID(
		id,
	)

	if err != nil {
		return err
	}

	if voter.IsSuspended {
		return errors.New(
			"suspended voter cannot enroll biometrics",
		)
	}

	voter.FaceEnrolled = true
	voter.UpdatedAt = time.Now()

	return s.VoterRepo.Update(
		*voter,
	)
}

// =====================================
// Verify Face
// =====================================

func (s *VoterService) VerifyFace(
	id int,
) error {

	if id <= 0 {
		return errors.New(
			"invalid voter ID",
		)
	}

	voter, err := s.VoterRepo.GetByID(
		id,
	)

	if err != nil {
		return err
	}

	if !voter.FaceEnrolled {
		return errors.New(
			"face has not been enrolled",
		)
	}

	if voter.IsSuspended {
		return errors.New(
			"suspended voter cannot be verified",
		)
	}

	voter.FaceVerified = true
	voter.UpdatedAt = time.Now()

	return s.VoterRepo.Update(
		*voter,
	)
}

// =====================================
// Enroll Fingerprint
// =====================================

func (s *VoterService) EnrollFingerprint(
	id int,
) error {

	if id <= 0 {
		return errors.New(
			"invalid voter ID",
		)
	}

	voter, err := s.VoterRepo.GetByID(
		id,
	)

	if err != nil {
		return err
	}

	if voter.IsSuspended {
		return errors.New(
			"suspended voter cannot enroll biometrics",
		)
	}

	voter.FingerprintEnrolled = true
	voter.UpdatedAt = time.Now()

	return s.VoterRepo.Update(
		*voter,
	)
}

// =====================================
// Verify Fingerprint
// =====================================

func (s *VoterService) VerifyFingerprint(
	id int,
) error {

	if id <= 0 {
		return errors.New(
			"invalid voter ID",
		)
	}

	voter, err := s.VoterRepo.GetByID(
		id,
	)

	if err != nil {
		return err
	}

	if !voter.FingerprintEnrolled {
		return errors.New(
			"fingerprint has not been enrolled",
		)
	}

	if voter.IsSuspended {
		return errors.New(
			"suspended voter cannot be verified",
		)
	}

	voter.FingerprintVerified = true
	voter.UpdatedAt = time.Now()

	return s.VoterRepo.Update(
		*voter,
	)
}

// =====================================
// Mark Voter As Voted
// =====================================

func (s *VoterService) MarkAsVoted(
	id int,
) error {

	if id <= 0 {
		return errors.New(
			"invalid voter ID",
		)
	}

	voter, err := s.VoterRepo.GetByID(
		id,
	)

	if err != nil {
		return err
	}

	if voter.IsSuspended {
		return errors.New(
			"suspended voter cannot vote",
		)
	}

	if !voter.IsEligible {
		return errors.New(
			"voter is not eligible",
		)
	}

	if voter.HasVoted {
		return errors.New(
			"voter has already voted",
		)
	}

	voter.HasVoted = true
	voter.LastVoteTime = time.Now()
	voter.UpdatedAt = time.Now()

	return s.VoterRepo.Update(
		*voter,
	)
}
