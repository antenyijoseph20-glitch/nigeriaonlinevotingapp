package services

import (
	"errors"
	"strings"
	"time"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/repositories"
)

type PositionService struct {
	PositionRepo repositories.PositionRepository
}

// =====================================
// Constructor
// =====================================

func NewPositionService(
	repo repositories.PositionRepository,
) *PositionService {

	return &PositionService{
		PositionRepo: repo,
	}
}

// =====================================
// Create Position
// =====================================

func (s *PositionService) CreatePosition(
	position models.Position,
) error {

	position.Name = strings.TrimSpace(
		position.Name,
	)

	position.Description = strings.TrimSpace(
		position.Description,
	)

	position.Level = strings.TrimSpace(
		position.Level,
	)

	if position.Name == "" {

		return errors.New(
			"position name is required",
		)
	}

	if position.Level == "" {

		return errors.New(
			"position level is required",
		)
	}

	if position.Seats <= 0 {

		return errors.New(
			"number of seats must be greater than zero",
		)
	}

	_, err := s.PositionRepo.GetByName(
		position.Name,
	)

	if err == nil {

		return errors.New(
			"position already exists",
		)
	}

	position.IsActive = true

	position.CreatedAt = time.Now()

	position.UpdatedAt = time.Now()

	return s.PositionRepo.Create(
		position,
	)
}

// =====================================
// Update Position
// =====================================

func (s *PositionService) UpdatePosition(
	position models.Position,
) error {

	position.Name = strings.TrimSpace(
		position.Name,
	)

	position.Description = strings.TrimSpace(
		position.Description,
	)

	position.Level = strings.TrimSpace(
		position.Level,
	)

	if position.Name == "" {

		return errors.New(
			"position name is required",
		)
	}

	if position.Level == "" {

		return errors.New(
			"position level is required",
		)
	}

	if position.Seats <= 0 {

		return errors.New(
			"number of seats must be greater than zero",
		)
	}

	position.UpdatedAt = time.Now()

	return s.PositionRepo.Update(
		position,
	)
}

// =====================================
// Activate Position
// =====================================

func (s *PositionService) ActivatePosition(
	id int,
) error {

	position, err := s.PositionRepo.GetByID(
		id,
	)

	if err != nil {

		return err
	}

	position.IsActive = true

	position.UpdatedAt = time.Now()

	return s.PositionRepo.Update(
		*position,
	)
}

// =====================================
// Deactivate Position
// =====================================

func (s *PositionService) DeactivatePosition(
	id int,
) error {

	position, err := s.PositionRepo.GetByID(
		id,
	)

	if err != nil {

		return err
	}

	position.IsActive = false

	position.UpdatedAt = time.Now()

	return s.PositionRepo.Update(
		*position,
	)
}

// =====================================
// Delete Position
// =====================================

func (s *PositionService) DeletePosition(
	id int,
) error {

	_, err := s.PositionRepo.GetByID(
		id,
	)

	if err != nil {

		return err
	}

	return s.PositionRepo.Delete(
		id,
	)
}

// =====================================
// Get Position By ID
// =====================================

func (s *PositionService) GetPositionByID(
	id int,
) (*models.Position, error) {

	return s.PositionRepo.GetByID(
		id,
	)
}

// =====================================
// Get Position By Name
// =====================================

func (s *PositionService) GetPositionByName(
	name string,
) (*models.Position, error) {

	return s.PositionRepo.GetByName(
		name,
	)
}

// =====================================
// Get All Positions
// =====================================

func (s *PositionService) GetAllPositions() []models.Position {

	return s.PositionRepo.GetAll()
}
