package services

import (
	"errors"
	"strings"
	"time"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/repositories"
)

type PartyService struct {
	PartyRepo repositories.PartyRepository
}

// =====================================
// Constructor
// =====================================

func NewPartyService(
	repo repositories.PartyRepository,
) *PartyService {

	return &PartyService{
		PartyRepo: repo,
	}
}

// =====================================
// Create Party
// =====================================

func (s *PartyService) CreateParty(
	party models.Party,
) error {

	party.Name = strings.TrimSpace(
		party.Name,
	)

	party.Abbreviation = strings.TrimSpace(
		party.Abbreviation,
	)

	party.Slogan = strings.TrimSpace(
		party.Slogan,
	)

	party.Logo = strings.TrimSpace(
		party.Logo,
	)

	party.Chairman = strings.TrimSpace(
		party.Chairman,
	)

	party.Headquarters = strings.TrimSpace(
		party.Headquarters,
	)

	party.Description = strings.TrimSpace(
		party.Description,
	)

	// -------------------------------------
	// Validate required fields
	// -------------------------------------

	if party.Name == "" {

		return errors.New(
			"party name is required",
		)
	}

	if party.Abbreviation == "" {

		return errors.New(
			"party abbreviation is required",
		)
	}

	// -------------------------------------
	// Check duplicate party name
	// -------------------------------------

	existing, err := s.PartyRepo.GetByName(
		party.Name,
	)

	if err == nil && existing != nil {

		return errors.New(
			"party name already exists",
		)
	}

	// -------------------------------------
	// Set default values
	// -------------------------------------

	party.IsActive = true

	party.CreatedAt = time.Now()

	party.UpdatedAt = time.Now()

	// -------------------------------------
	// Save party
	// -------------------------------------

	return s.PartyRepo.Create(
		party,
	)
}

// =====================================
// Update Party
// =====================================

func (s *PartyService) UpdateParty(
	party models.Party,
) error {

	party.Name = strings.TrimSpace(
		party.Name,
	)

	party.Abbreviation = strings.TrimSpace(
		party.Abbreviation,
	)

	party.Slogan = strings.TrimSpace(
		party.Slogan,
	)

	party.Logo = strings.TrimSpace(
		party.Logo,
	)

	party.Chairman = strings.TrimSpace(
		party.Chairman,
	)

	party.Headquarters = strings.TrimSpace(
		party.Headquarters,
	)

	party.Description = strings.TrimSpace(
		party.Description,
	)

	// -------------------------------------
	// Validate ID
	// -------------------------------------

	if party.ID <= 0 {

		return errors.New(
			"invalid party ID",
		)
	}

	// -------------------------------------
	// Validate required fields
	// -------------------------------------

	if party.Name == "" {

		return errors.New(
			"party name is required",
		)
	}

	if party.Abbreviation == "" {

		return errors.New(
			"party abbreviation is required",
		)
	}

	// -------------------------------------
	// Verify party exists
	// -------------------------------------

	existing, err := s.PartyRepo.GetByID(
		party.ID,
	)

	if err != nil {

		return err
	}

	// -------------------------------------
	// Preserve original creation time
	// -------------------------------------

	party.CreatedAt = existing.CreatedAt

	party.UpdatedAt = time.Now()

	// -------------------------------------
	// Update party
	// -------------------------------------

	return s.PartyRepo.Update(
		party,
	)
}

// =====================================
// Delete Party
// =====================================

func (s *PartyService) DeleteParty(
	id int,
) error {

	if id <= 0 {

		return errors.New(
			"invalid party ID",
		)
	}

	_, err := s.PartyRepo.GetByID(
		id,
	)

	if err != nil {

		return err
	}

	return s.PartyRepo.Delete(
		id,
	)
}

// =====================================
// Get Party By ID
// =====================================

func (s *PartyService) GetPartyByID(
	id int,
) (*models.Party, error) {

	return s.PartyRepo.GetByID(
		id,
	)
}

// =====================================
// Get All Parties
// =====================================

func (s *PartyService) GetAllParties() []models.Party {

	return s.PartyRepo.GetAll()
}

// =====================================
// Activate Party
// =====================================

func (s *PartyService) ActivateParty(
	id int,
) error {

	party, err := s.PartyRepo.GetByID(
		id,
	)

	if err != nil {

		return err
	}

	party.IsActive = true

	party.UpdatedAt = time.Now()

	return s.PartyRepo.Update(
		*party,
	)
}

// =====================================
// Deactivate Party
// =====================================

func (s *PartyService) DeactivateParty(
	id int,
) error {

	party, err := s.PartyRepo.GetByID(
		id,
	)

	if err != nil {

		return err
	}

	party.IsActive = false

	party.UpdatedAt = time.Now()

	return s.PartyRepo.Update(
		*party,
	)
}
