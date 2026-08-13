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

	// -------------------------------------
	// Normalize input
	// -------------------------------------

	position.Name = strings.TrimSpace(position.Name)
	position.Description = strings.TrimSpace(position.Description)
	position.Level = strings.TrimSpace(position.Level)

	// -------------------------------------
	// Validate required fields
	// -------------------------------------

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

	// -------------------------------------
	// Validate level
	// -------------------------------------

	if !isValidPositionLevel(position.Level) {
		return errors.New(
			"invalid position level",
		)
	}

	// -------------------------------------
	// Check duplicate position
	// -------------------------------------

	existing, err := s.PositionRepo.GetByName(
		position.Name,
	)

	if err == nil && existing != nil {
		return errors.New(
			"position already exists",
		)
	}

	// A repository error other than "not found"
	// should not be silently ignored.
	if err != nil &&
		err.Error() != "position not found" {

		return err
	}

	// -------------------------------------
	// Set system-controlled fields
	// -------------------------------------

	position.IsActive = true

	now := time.Now()

	position.CreatedAt = now
	position.UpdatedAt = now

	// -------------------------------------
	// Persist
	// -------------------------------------

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

	// -------------------------------------
	// Validate ID
	// -------------------------------------

	if position.ID <= 0 {
		return errors.New(
			"invalid position ID",
		)
	}

	// -------------------------------------
	// Normalize input
	// -------------------------------------

	position.Name = strings.TrimSpace(position.Name)
	position.Description = strings.TrimSpace(position.Description)
	position.Level = strings.TrimSpace(position.Level)

	// -------------------------------------
	// Validate required fields
	// -------------------------------------

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

	if !isValidPositionLevel(position.Level) {
		return errors.New(
			"invalid position level",
		)
	}

	// -------------------------------------
	// Get existing position
	// -------------------------------------

	existing, err := s.PositionRepo.GetByID(
		position.ID,
	)

	if err != nil {
		return err
	}

	// -------------------------------------
	// Prevent duplicate names
	// -------------------------------------

	other, err := s.PositionRepo.GetByName(
		position.Name,
	)

	if err == nil &&
		other != nil &&
		other.ID != position.ID {

		return errors.New(
			"position name already exists",
		)
	}

	if err != nil &&
		err.Error() != "position not found" {

		return err
	}

	// -------------------------------------
	// Preserve system-controlled fields
	// -------------------------------------

	position.CreatedAt = existing.CreatedAt

	position.UpdatedAt = time.Now()

	// -------------------------------------
	// Persist
	// -------------------------------------

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

	if id <= 0 {
		return errors.New(
			"invalid position ID",
		)
	}

	position, err := s.PositionRepo.GetByID(
		id,
	)

	if err != nil {
		return err
	}

	// Already active is not an error.
	if position.IsActive {
		return nil
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

	if id <= 0 {
		return errors.New(
			"invalid position ID",
		)
	}

	position, err := s.PositionRepo.GetByID(
		id,
	)

	if err != nil {
		return err
	}

	// Already inactive is not an error.
	if !position.IsActive {
		return nil
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

	if id <= 0 {
		return errors.New(
			"invalid position ID",
		)
	}

	position, err := s.PositionRepo.GetByID(
		id,
	)

	if err != nil {
		return err
	}

	// Active positions should be deactivated
	// instead of being deleted.
	if position.IsActive {
		return errors.New(
			"cannot delete an active position",
		)
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

	if id <= 0 {
		return nil, errors.New(
			"invalid position ID",
		)
	}

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

	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errors.New(
			"position name is required",
		)
	}

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

// =====================================
// Position Level Validation
// =====================================

func isValidPositionLevel(
	level string,
) bool {

	switch strings.ToLower(
		strings.TrimSpace(level),
	) {

	case "federal":
		return true

	case "state":
		return true

	case "local":
		return true

	case "local government":
		return true

	default:
		return false
	}
}
