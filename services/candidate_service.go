package services

import (
	"errors"
	"strings"
	"time"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/repositories"
)

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

	// Validate Election
	if _, err := s.ElectionRepo.GetByID(candidate.ElectionID); err != nil {
		return errors.New("selected election does not exist")
	}

	// Validate Political Party
	if _, err := s.PartyRepo.GetByID(candidate.PartyID); err != nil {
		return errors.New("selected political party does not exist")
	}

	// Validate Position
	if _, err := s.PositionRepo.GetByID(candidate.PositionID); err != nil {
		return errors.New("selected position does not exist")
	}

	candidate.IsApproved = false
	candidate.IsActive = true

	candidate.CreatedAt = time.Now()
	candidate.UpdatedAt = time.Now()

	return s.CandidateRepo.Create(candidate)
}

// =====================================
// Update Candidate
// =====================================

func (s *CandidateService) UpdateCandidate(
	candidate models.Candidate,
) error {

	candidate.UpdatedAt = time.Now()

	return s.CandidateRepo.Update(candidate)
}

// =====================================
// Delete Candidate
// =====================================

func (s *CandidateService) DeleteCandidate(
	id int,
) error {

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

	candidate, err := s.CandidateRepo.GetByID(id)

	if err != nil {
		return err
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

	candidate, err := s.CandidateRepo.GetByID(id)

	if err != nil {
		return err
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

	candidate, err := s.CandidateRepo.GetByID(id)

	if err != nil {
		return err
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

	return s.CandidateRepo.GetByID(id)
}

// =====================================
// Get All Candidates
// =====================================

func (s *CandidateService) GetAllCandidates() []models.Candidate {

	return s.CandidateRepo.GetAll()
}

// =====================================
// Get Candidates By Election
// =====================================

func (s *CandidateService) GetCandidatesByElection(
	electionID int,
) []models.Candidate {

	return s.CandidateRepo.GetByElectionID(electionID)
}

// =====================================
// Get Candidates By Party
// =====================================

func (s *CandidateService) GetCandidatesByParty(
	partyID int,
) []models.Candidate {

	return s.CandidateRepo.GetByPartyID(partyID)
}

// =====================================
// Get Candidates By Position
// =====================================

func (s *CandidateService) GetCandidatesByPosition(
	positionID int,
) []models.Candidate {

	return s.CandidateRepo.GetByPositionID(positionID)
}

// =====================================
// Get All Candidate Views
// =====================================

func (s *CandidateService) GetAllCandidateViews() []models.CandidateView {

	candidates := s.CandidateRepo.GetAll()

	views := []models.CandidateView{}

	for _, candidate := range candidates {

		view := models.CandidateView{

			ID: candidate.ID,

			FirstName: candidate.FirstName,

			LastName: candidate.LastName,

			Gender: candidate.Gender,

			Email: candidate.Email,

			PhoneNumber: candidate.PhoneNumber,

			IsApproved: candidate.IsApproved,

			IsActive: candidate.IsActive,
		}

		// -------------------------
		// Election
		// -------------------------

		election, err := s.ElectionRepo.GetByID(
			candidate.ElectionID,
		)

		if err == nil {

			view.ElectionName = election.Title

		} else {

			view.ElectionName = "Unknown Election"

		}

		// -------------------------
		// Political Party
		// -------------------------

		party, err := s.PartyRepo.GetByID(
			candidate.PartyID,
		)

		if err == nil {

			view.PartyName = party.Name

		} else {

			view.PartyName = "Unknown Party"

		}

		// -------------------------
		// Position
		// -------------------------

		position, err := s.PositionRepo.GetByID(
			candidate.PositionID,
		)

		if err == nil {

			view.PositionName = position.Name

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
