package repositories

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"sync"

	"nigeriaonlinevoting/models"
)

// PositionRepository defines operations
// for managing elective positions.
type PositionRepository interface {

	// Create a position.
	Create(
		position models.Position,
	) error

	// Update a position.
	Update(
		position models.Position,
	) error

	// Delete a position.
	Delete(
		id int,
	) error

	// Get position by ID.
	GetByID(
		id int,
	) (*models.Position, error)

	// Get position by name.
	GetByName(
		name string,
	) (*models.Position, error)

	// Get all positions.
	GetAll() []models.Position
}

// =====================================
// JSON Repository
// =====================================

type PositionJSONRepository struct {
	filePath string
	mutex    sync.Mutex
}

// =====================================
// Constructor
// =====================================

func NewPositionRepository(
	filePath string,
) *PositionJSONRepository {

	return &PositionJSONRepository{
		filePath: filePath,
	}
}

// =====================================
// Load
// =====================================

func (r *PositionJSONRepository) load() ([]models.Position, error) {

	file, err := os.Open(r.filePath)

	if os.IsNotExist(err) {
		return []models.Position{}, nil
	}

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var positions []models.Position

	err = json.NewDecoder(file).Decode(&positions)

	if errors.Is(err, io.EOF) {
		return []models.Position{}, nil
	}

	if err != nil {
		return nil, err
	}

	if positions == nil {
		return []models.Position{}, nil
	}

	return positions, nil
}

// =====================================
// Save
// =====================================

func (r *PositionJSONRepository) save(
	positions []models.Position,
) error {

	file, err := os.Create(r.filePath)

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	encoder.SetIndent("", "    ")

	return encoder.Encode(positions)
}

// =====================================
// Create
// =====================================

func (r *PositionJSONRepository) Create(
	position models.Position,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Normalize input.
	position.Name = strings.TrimSpace(position.Name)
	position.Description = strings.TrimSpace(position.Description)
	position.Level = strings.TrimSpace(position.Level)

	// -------------------------
	// Basic validation
	// -------------------------

	if position.Name == "" {
		return errors.New("position name is required")
	}

	if position.Seats <= 0 {
		return errors.New("position seats must be greater than zero")
	}

	if position.Level == "" {
		return errors.New("position level is required")
	}

	positions, err := r.load()

	if err != nil {
		return err
	}

	// -------------------------
	// Prevent duplicate names
	// -------------------------

	normalizedName := strings.ToLower(position.Name)

	for _, existing := range positions {

		if strings.ToLower(
			strings.TrimSpace(existing.Name),
		) == normalizedName {

			return errors.New(
				"position already exists",
			)
		}
	}

	// -------------------------
	// Generate safe ID
	// -------------------------

	maxID := 0

	for _, existing := range positions {

		if existing.ID > maxID {
			maxID = existing.ID
		}
	}

	position.ID = maxID + 1

	positions = append(
		positions,
		position,
	)

	return r.save(positions)
}

// =====================================
// Update
// =====================================

func (r *PositionJSONRepository) Update(
	position models.Position,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	position.Name = strings.TrimSpace(position.Name)
	position.Description = strings.TrimSpace(position.Description)
	position.Level = strings.TrimSpace(position.Level)

	// -------------------------
	// Basic validation
	// -------------------------

	if position.ID <= 0 {
		return errors.New("invalid position ID")
	}

	if position.Name == "" {
		return errors.New("position name is required")
	}

	if position.Seats <= 0 {
		return errors.New("position seats must be greater than zero")
	}

	if position.Level == "" {
		return errors.New("position level is required")
	}

	positions, err := r.load()

	if err != nil {
		return err
	}

	targetIndex := -1

	for i := range positions {

		if positions[i].ID == position.ID {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		return errors.New("position not found")
	}

	// -------------------------
	// Prevent duplicate names
	// -------------------------

	normalizedName := strings.ToLower(position.Name)

	for i, existing := range positions {

		if i == targetIndex {
			continue
		}

		existingName := strings.ToLower(
			strings.TrimSpace(existing.Name),
		)

		if existingName == normalizedName {

			return errors.New(
				"position name already exists",
			)
		}
	}

	positions[targetIndex] = position

	return r.save(positions)
}

// =====================================
// Delete
// =====================================

func (r *PositionJSONRepository) Delete(
	id int,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if id <= 0 {
		return errors.New("invalid position ID")
	}

	positions, err := r.load()

	if err != nil {
		return err
	}

	for i := range positions {

		if positions[i].ID == id {

			positions = append(
				positions[:i],
				positions[i+1:]...,
			)

			return r.save(positions)
		}
	}

	return errors.New("position not found")
}

// =====================================
// Get By ID
// =====================================

func (r *PositionJSONRepository) GetByID(
	id int,
) (*models.Position, error) {

	if id <= 0 {
		return nil, errors.New(
			"invalid position ID",
		)
	}

	positions, err := r.load()

	if err != nil {
		return nil, err
	}

	for i := range positions {

		if positions[i].ID == id {
			return &positions[i], nil
		}
	}

	return nil, errors.New(
		"position not found",
	)
}

// =====================================
// Get By Name
// =====================================

func (r *PositionJSONRepository) GetByName(
	name string,
) (*models.Position, error) {

	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errors.New(
			"position name is required",
		)
	}

	positions, err := r.load()

	if err != nil {
		return nil, err
	}

	normalizedName := strings.ToLower(name)

	for i := range positions {

		existingName := strings.ToLower(
			strings.TrimSpace(positions[i].Name),
		)

		if existingName == normalizedName {
			return &positions[i], nil
		}
	}

	return nil, errors.New(
		"position not found",
	)
}

// =====================================
// Get All
// =====================================

func (r *PositionJSONRepository) GetAll() []models.Position {

	positions, err := r.load()

	if err != nil {
		return []models.Position{}
	}

	return positions
}
