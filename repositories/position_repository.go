package repositories

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"nigeriaonlinevoting/models"
)

type PositionRepository interface {

	// Create
	Create(
		position models.Position,
	) error

	// Update
	Update(
		position models.Position,
	) error

	// Delete
	Delete(
		id int,
	) error

	// Get by ID
	GetByID(
		id int,
	) (*models.Position, error)

	// Get by Name
	GetByName(
		name string,
	) (*models.Position, error)

	// Get All
	GetAll() []models.Position
}

type PositionJSONRepository struct {
	filePath string
	mutex    sync.Mutex
}

// Constructor

func NewPositionRepository(
	filePath string,
) *PositionJSONRepository {

	return &PositionJSONRepository{
		filePath: filePath,
	}
}

// ==========================
// Load
// ==========================

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

	err = json.NewDecoder(file).Decode(
		&positions,
	)

	if err != nil {

		return []models.Position{}, nil
	}

	return positions, nil
}

// ==========================
// Save
// ==========================

func (r *PositionJSONRepository) save(
	positions []models.Position,
) error {

	file, err := os.Create(
		r.filePath,
	)

	if err != nil {

		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(
		file,
	)

	encoder.SetIndent(
		"",
		"    ",
	)

	return encoder.Encode(
		positions,
	)
}

// ==========================
// Create
// ==========================

func (r *PositionJSONRepository) Create(
	position models.Position,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	positions, err := r.load()

	if err != nil {

		return err
	}

	position.ID = len(positions) + 1

	positions = append(
		positions,
		position,
	)

	return r.save(
		positions,
	)
}

// ==========================
// Update
// ==========================

func (r *PositionJSONRepository) Update(
	position models.Position,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	positions, err := r.load()

	if err != nil {

		return err
	}

	for i := range positions {

		if positions[i].ID == position.ID {

			positions[i] = position

			return r.save(
				positions,
			)
		}
	}

	return errors.New(
		"position not found",
	)
}

// ==========================
// Delete
// ==========================

func (r *PositionJSONRepository) Delete(
	id int,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

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

			return r.save(
				positions,
			)
		}
	}

	return errors.New(
		"position not found",
	)
}

// ==========================
// Get By ID
// ==========================

func (r *PositionJSONRepository) GetByID(
	id int,
) (*models.Position, error) {

	positions, err := r.load()

	if err != nil {

		return nil, err
	}

	for _, position := range positions {

		if position.ID == id {

			return &position, nil
		}
	}

	return nil, errors.New(
		"position not found",
	)
}

// ==========================
// Get By Name
// ==========================

func (r *PositionJSONRepository) GetByName(
	name string,
) (*models.Position, error) {

	positions, err := r.load()

	if err != nil {

		return nil, err
	}

	for _, position := range positions {

		if position.Name == name {

			return &position, nil
		}
	}

	return nil, errors.New(
		"position not found",
	)
}

// ==========================
// Get All
// ==========================

func (r *PositionJSONRepository) GetAll() []models.Position {

	positions, err := r.load()

	if err != nil {

		return []models.Position{}
	}

	return positions
}
