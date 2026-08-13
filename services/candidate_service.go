package services

import (
	"errors"
	"strings"
	"time"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/repositories"
)

// CandidateService handles candidate business rules.
type CandidateService struct {
	CandidateRepo repositories.CandidateRepository
	ElectionRepo  repositories.ElectionRepository
	PartyRepo     repositories.PartyRepository
	PositionRepo  repositories.PositionRepository
}

// =====================================
// Constructor
// =====================================

func NewCandidateService(
	candidateRepo repositories.CandidateRepository,
	electionRepo repositories.ElectionRepository,
	partyRepo repositories.PartyRepository,
	positionRepo repositories.PositionRepository,
) *CandidateService {

	return &CandidateService{
		CandidateRepo: candidateRepo,
		ElectionRepo:  electionRepo,
		PartyRepo:     partyRepo,
		PositionRepo:  positionRepo,
	}
}

// =====================================
// Create Candidate
// =====================================

func (s *CandidateService) CreateCandidate(
	candidate models.Candidate,
) error {

	// -------------------------------------
	// Basic validation
	// -------------------------------------

	candidate.FirstName = strings.TrimSpace(candidate.FirstName)
	candidate.LastName = strings.TrimSpace(candidate.LastName)
	candidate.Email = strings.TrimSpace(candidate.Email)
	candidate.PhoneNumber = strings.TrimSpace(candidate.PhoneNumber)
	candidate.Biography = strings.TrimSpace(candidate.Biography)
	candidate.Manifesto = strings.TrimSpace(candidate.Manifesto)

	if candidate.FirstName == "" {
		return errors.New("first name is required")
	}

	if candidate.LastName == "" {
		return errors.New("last name is required")
	}

	if candidate.ElectionID <= 0 {
		return errors.New("invalid election")
	}

	if candidate.PartyID <= 0 {
		return errors.New("invalid political party")
	}

	if candidate.PositionID <= 0 {
		return errors.New("invalid position")
	}

	// -------------------------------------
	// Validate dependencies
	// -------------------------------------

	election, err := s.ElectionRepo.GetByID(
		candidate.ElectionID,
	)

	if err != nil || election == nil {
		return errors.New(
			"selected election does not exist",
		)
	}

	party, err := s.PartyRepo.GetByID(
		candidate.PartyID,
	)

	if err != nil || party == nil {
		return errors.New(
			"selected political party does not exist",
		)
	}

	position, err := s.PositionRepo.GetByID(
		candidate.PositionID,
	)

	if err != nil || position == nil {
		return errors.New(
			"selected position does not exist",
		)
	}

	// -------------------------------------
	// Election state
	// -------------------------------------

	if election.Status == "closed" {
		return errors.New(
			"cannot add candidate to a closed election",
		)
	}

	// -------------------------------------
	// Position state
	// -------------------------------------

	if !position.IsActive {
		return errors.New(
			"selected position is inactive",
		)
	}

	// -------------------------------------
	// Candidate defaults
	// -------------------------------------

	candidate.IsApproved = false
	candidate.IsActive = true

	now := time.Now()

	candidate.CreatedAt = now
	candidate.UpdatedAt = now

	// -------------------------------------
	// Store candidate
	// -------------------------------------

	return s.CandidateRepo.Create(candidate)
}

// =====================================
// Update Candidate
// =====================================

func (s *CandidateService) UpdateCandidate(
	candidate models.Candidate,
) error {

	// -------------------------------------
	// Validate ID
	// -------------------------------------

	if candidate.ID <= 0 {
		return errors.New("invalid candidate ID")
	}

	// -------------------------------------
	// Clean strings
	// -------------------------------------

	candidate.FirstName = strings.TrimSpace(candidate.FirstName)
	candidate.LastName = strings.TrimSpace(candidate.LastName)
	candidate.Email = strings.TrimSpace(candidate.Email)
	candidate.PhoneNumber = strings.TrimSpace(candidate.PhoneNumber)
	candidate.Biography = strings.TrimSpace(candidate.Biography)
	candidate.Manifesto = strings.TrimSpace(candidate.Manifesto)

	if candidate.FirstName == "" {
		return errors.New("first name is required")
	}

	if candidate.LastName == "" {
		return errors.New("last name is required")
	}

	// -------------------------------------
	// Confirm candidate exists
	// -------------------------------------

	existing, err := s.CandidateRepo.GetByID(
		candidate.ID,
	)

	if err != nil || existing == nil {
		return errors.New(
			"candidate not found",
		)
	}

	// -------------------------------------
	// Validate election
	// -------------------------------------

	if candidate.ElectionID <= 0 {
		return errors.New("invalid election")
	}

	election, err := s.ElectionRepo.GetByID(
		candidate.ElectionID,
	)

	if err != nil || election == nil {
		return errors.New(
			"selected election does not exist",
		)
	}

	if election.Status == "closed" {
		return errors.New(
			"cannot update candidate in a closed election",
		)
	}

	// -------------------------------------
	// Validate party
	// -------------------------------------

	if candidate.PartyID <= 0 {
		return errors.New("invalid political party")
	}

	party, err := s.PartyRepo.GetByID(
		candidate.PartyID,
	)

	if err != nil || party == nil {
		return errors.New(
			"selected political party does not exist",
		)
	}

	// -------------------------------------
	// Validate position
	// -------------------------------------

	if candidate.PositionID <= 0 {
		return errors.New("invalid position")
	}

	position, err := s.PositionRepo.GetByID(
		candidate.PositionID,
	)

	if err != nil || position == nil {
		return errors.New(
			"selected position does not exist",
		)
	}

	if !position.IsActive {
		return errors.New(
			"selected position is inactive",
		)
	}

	// -------------------------------------
	// Preserve approval status
	// -------------------------------------

	candidate.IsApproved = existing.IsApproved

	// -------------------------------------
	// Preserve active status
	// -------------------------------------

	candidate.IsActive = existing.IsActive

	// -------------------------------------
	// Preserve creation time
	// -------------------------------------

	candidate.CreatedAt = existing.CreatedAt

	if candidate.CreatedAt.IsZero() {
		candidate.CreatedAt = time.Now()
	}

	candidate.UpdatedAt = time.Now()

	// -------------------------------------
	// Save
	// -------------------------------------

	return s.CandidateRepo.Update(candidate)
}

// =====================================
// Delete Candidate
// =====================================

func (s *CandidateService) DeleteCandidate(
	id int,
) error {

	if id <= 0 {
		return errors.New(
			"invalid candidate ID",
		)
	}

	_, err := s.CandidateRepo.GetByID(id)

	if err != nil {
		return err
	}

	return s.CandidateRepo.Delete(id)
}

// =====================================
// Approve Candidate
// =====================================

func (s *CandidateService) ApproveCandidate(
	id int,
) error {

	if id <= 0 {
		return errors.New(
			"invalid candidate ID",
		)
	}

	candidate, err := s.CandidateRepo.GetByID(id)

	if err != nil || candidate == nil {
		return errors.New(
			"candidate not found",
		)
	}

	// Already approved is not a crash,
	// but there is no need to update again.
	if candidate.IsApproved {
		return errors.New(
			"candidate is already approved",
		)
	}

	// Candidate must belong to a valid election.
	election, err := s.ElectionRepo.GetByID(
		candidate.ElectionID,
	)

	if err != nil || election == nil {
		return errors.New(
			"candidate election does not exist",
		)
	}

	// Candidate must belong to a valid position.
	position, err := s.PositionRepo.GetByID(
		candidate.PositionID,
	)

	if err != nil || position == nil {
		return errors.New(
			"candidate position does not exist",
		)
	}

	if !position.IsActive {
		return errors.New(
			"candidate position is inactive",
		)
	}

	candidate.IsApproved = true
	candidate.UpdatedAt = time.Now()

	return s.CandidateRepo.Update(*candidate)
}

// =====================================
// Deactivate Candidate
// =====================================

func (s *CandidateService) DeactivateCandidate(
	id int,
) error {

	if id <= 0 {
		return errors.New(
			"invalid candidate ID",
		)
	}

	candidate, err := s.CandidateRepo.GetByID(id)

	if err != nil || candidate == nil {
		return errors.New(
			"candidate not found",
		)
	}

	if !candidate.IsActive {
		return errors.New(
			"candidate is already inactive",
		)
	}

	candidate.IsActive = false
	candidate.UpdatedAt = time.Now()

	return s.CandidateRepo.Update(*candidate)
}

// =====================================
// Activate Candidate
// =====================================

func (s *CandidateService) ActivateCandidate(
	id int,
) error {

	if id <= 0 {
		return errors.New(
			"invalid candidate ID",
		)
	}

	candidate, err := s.CandidateRepo.GetByID(id)

	if err != nil || candidate == nil {
		return errors.New(
			"candidate not found",
		)
	}

	if candidate.IsActive {
		return errors.New(
			"candidate is already active",
		)
	}

	// Do not activate an unapproved candidate.
	if !candidate.IsApproved {
		return errors.New(
			"candidate must be approved before activation",
		)
	}

	candidate.IsActive = true
	candidate.UpdatedAt = time.Now()

	return s.CandidateRepo.Update(*candidate)
}

// =====================================
// Get Candidate By ID
// =====================================

func (s *CandidateService) GetCandidateByID(
	id int,
) (*models.Candidate, error) {

	if id <= 0 {
		return nil, errors.New(
			"invalid candidate ID",
		)
	}

	return s.CandidateRepo.GetByID(id)
}

// =====================================
// Get All Candidates
// =====================================

func (s *CandidateService) GetAllCandidates() []models.Candidate {

	candidates := s.CandidateRepo.GetAll()

	if candidates == nil {
		return []models.Candidate{}
	}

	return candidates
}

// =====================================
// Get Candidates By Election
// =====================================

func (s *CandidateService) GetCandidatesByElection(
	electionID int,
) []models.Candidate {

	if electionID <= 0 {
		return []models.Candidate{}
	}

	candidates := s.CandidateRepo.GetByElectionID(
		electionID,
	)

	if candidates == nil {
		return []models.Candidate{}
	}

	return candidates
}

// =====================================
// Get Candidates By Party
// =====================================

func (s *CandidateService) GetCandidatesByParty(
	partyID int,
) []models.Candidate {

	if partyID <= 0 {
		return []models.Candidate{}
	}

	candidates := s.CandidateRepo.GetByPartyID(
		partyID,
	)

	if candidates == nil {
		return []models.Candidate{}
	}

	return candidates
}

// =====================================
// Get Candidates By Position
// =====================================

func (s *CandidateService) GetCandidatesByPosition(
	positionID int,
) []models.Candidate {

	if positionID <= 0 {
		return []models.Candidate{}
	}

	candidates := s.CandidateRepo.GetByPositionID(
		positionID,
	)

	if candidates == nil {
		return []models.Candidate{}
	}

	return candidates
}

// =====================================
// Get All Candidate Views
// =====================================

func (s *CandidateService) GetAllCandidateViews() []models.CandidateView {

	candidates := s.CandidateRepo.GetAll()

	views := make(
		[]models.CandidateView,
		0,
		len(candidates),
	)

	for _, candidate := range candidates {

		view := models.CandidateView{
			ID:          candidate.ID,
			FirstName:   candidate.FirstName,
			LastName:    candidate.LastName,
			Gender:      candidate.Gender,
			Email:       candidate.Email,
			PhoneNumber: candidate.PhoneNumber,
			IsApproved:  candidate.IsApproved,
			IsActive:    candidate.IsActive,
		}

		// -------------------------------------
		// Election
		// -------------------------------------

		if candidate.ElectionID > 0 {

			election, err := s.ElectionRepo.GetByID(
				candidate.ElectionID,
			)

			if err == nil && election != nil {
				view.ElectionName = election.Title
			} else {
				view.ElectionName = "Unknown Election"
			}

		} else {
			view.ElectionName = "Unknown Election"
		}

		// -------------------------------------
		// Political Party
		// -------------------------------------

		if candidate.PartyID > 0 {

			party, err := s.PartyRepo.GetByID(
				candidate.PartyID,
			)

			if err == nil && party != nil {
				view.PartyName = party.Name
			} else {
				view.PartyName = "Unknown Party"
			}

		} else {
			view.PartyName = "Unknown Party"
		}

		// -------------------------------------
		// Position
		// -------------------------------------

		if candidate.PositionID > 0 {

			position, err := s.PositionRepo.GetByID(
				candidate.PositionID,
			)

			if err == nil && position != nil {
				view.PositionName = position.Name
			} else {
				view.PositionName = "Unknown Position"
			}

		} else {
			view.PositionName = "Unknown Position"
		}

		views = append(
			views,
			view,
		)
	}

	return views
}
