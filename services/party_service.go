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

// Constructor

func NewPartyService(
	repo repositories.PartyRepository,
) *PartyService {

	return &PartyService{
		PartyRepo: repo,
	}
}

// CreateParty

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

	// Check duplicate name

	existing, _ := s.PartyRepo.GetByName(
		party.Name,
	)

	if existing != nil {

		return errors.New(
			"party already exists",
		)
	}

	party.IsActive = true

	party.CreatedAt = time.Now()

	party.UpdatedAt = time.Now()

	return s.PartyRepo.Create(
		party,
	)
}

// UpdateParty

func (s *PartyService) UpdateParty(
	party models.Party,
) error {

	party.UpdatedAt = time.Now()

	return s.PartyRepo.Update(
		party,
	)
}

// DeleteParty

func (s *PartyService) DeleteParty(
	id int,
) error {

	return s.PartyRepo.Delete(
		id,
	)
}

// GetPartyByID

func (s *PartyService) GetPartyByID(
	id int,
) (*models.Party, error) {

	return s.PartyRepo.GetByID(
		id,
	)
}

// GetAllParties

func (s *PartyService) GetAllParties() []models.Party {

	return s.PartyRepo.GetAll()
}

// ActivateParty

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

// DeactivateParty

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
