package services

import (
	"errors"
	"strings"
	"time"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/repositories"
)

type ElectionService struct {
	ElectionRepo repositories.ElectionRepository
}

// Constructor

func NewElectionService(
	repo repositories.ElectionRepository,
) *ElectionService {

	return &ElectionService{
		ElectionRepo: repo,
	}
}

// CreateElection creates a new election.

func (s *ElectionService) CreateElection(
	election models.Election,
) error {

	election.Title = strings.TrimSpace(
		election.Title,
	)

	election.Description = strings.TrimSpace(
		election.Description,
	)

	if election.Title == "" {

		return errors.New(
			"election title is required",
		)
	}

	if election.StartDate.After(
		election.EndDate,
	) {

		return errors.New(
			"end date must be after start date",
		)
	}

	election.Status = "draft"

	election.CreatedAt = time.Now()

	election.UpdatedAt = time.Now()

	return s.ElectionRepo.Create(
		election,
	)
}

// GetElectionByID

func (s *ElectionService) GetElectionByID(
	id int,
) (*models.Election, error) {

	return s.ElectionRepo.GetByID(
		id,
	)
}

// GetCurrentElection

func (s *ElectionService) GetCurrentElection() (*models.Election, error) {

	return s.ElectionRepo.GetCurrent()
}

// GetAllElections

func (s *ElectionService) GetAllElections() []models.Election {

	return s.ElectionRepo.GetAll()
}

// UpdateElection

func (s *ElectionService) UpdateElection(
	election models.Election,
) error {

	election.UpdatedAt = time.Now()

	return s.ElectionRepo.Update(
		election,
	)
}

// DeleteElection

func (s *ElectionService) DeleteElection(
	id int,
) error {

	election, err := s.ElectionRepo.GetByID(
		id,
	)

	if err != nil {

		return err
	}

	if election.Status == "open" {

		return errors.New(
			"cannot delete an open election",
		)
	}

	return s.ElectionRepo.Delete(
		id,
	)
}

// OpenElection

func (s *ElectionService) OpenElection(
	id int,
) error {

	current, _ := s.ElectionRepo.GetCurrent()

	if current != nil {

		return errors.New(
			"another election is already open",
		)
	}

	election, err := s.ElectionRepo.GetByID(
		id,
	)

	if err != nil {

		return err
	}

	election.Status = "open"

	election.UpdatedAt = time.Now()

	return s.ElectionRepo.Update(
		*election,
	)
}

// CloseElection

func (s *ElectionService) CloseElection(
	id int,
) error {

	election, err := s.ElectionRepo.GetByID(
		id,
	)

	if err != nil {

		return err
	}

	election.Status = "closed"

	election.UpdatedAt = time.Now()

	return s.ElectionRepo.Update(
		*election,
	)
}
