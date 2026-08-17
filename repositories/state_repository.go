package repositories

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"

	"nigeriaonlinevoting/models"
)

// StateRepository defines operations for reading
// and managing Nigeria's states and FCT.
type StateRepository interface {
	GetByID(id int) (*models.State, error)
	GetByCode(code string) (*models.State, error)
	GetByName(name string) (*models.State, error)
	GetAll() []models.State
}

// StateJSONRepository stores state data in JSON.
type StateJSONRepository struct {
	filePath string
	mutex    sync.RWMutex
}

// NewStateRepository creates a new JSON-backed
// state repository.
func NewStateRepository(filePath string) *StateJSONRepository {
	return &StateJSONRepository{
		filePath: filePath,
	}
}

// load reads states from the JSON file.
func (r *StateJSONRepository) load() ([]models.State, error) {
	file, err := os.Open(r.filePath)

	if os.IsNotExist(err) {
		return []models.State{}, nil
	}

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var states []models.State

	err = json.NewDecoder(file).Decode(&states)

	if err != nil {
		return nil, err
	}

	if states == nil {
		return []models.State{}, nil
	}

	return states, nil
}

// GetByID returns a state by its internal ID.
func (r *StateJSONRepository) GetByID(
	id int,
) (*models.State, error) {

	if id <= 0 {
		return nil, errors.New("invalid state ID")
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	states, err := r.load()

	if err != nil {
		return nil, err
	}

	for i := range states {
		if states[i].ID == id {
			state := states[i]
			return &state, nil
		}
	}

	return nil, errors.New("state not found")
}

// GetByCode returns a state by its code.
func (r *StateJSONRepository) GetByCode(
	code string,
) (*models.State, error) {

	code = strings.TrimSpace(code)

	if code == "" {
		return nil, errors.New("state code is required")
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	states, err := r.load()

	if err != nil {
		return nil, err
	}

	for i := range states {
		if strings.EqualFold(
			strings.TrimSpace(states[i].Code),
			code,
		) {
			state := states[i]
			return &state, nil
		}
	}

	return nil, errors.New("state not found")
}

// GetByName returns a state by its name.
func (r *StateJSONRepository) GetByName(
	name string,
) (*models.State, error) {

	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errors.New("state name is required")
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	states, err := r.load()

	if err != nil {
		return nil, err
	}

	for i := range states {
		if strings.EqualFold(
			strings.TrimSpace(states[i].Name),
			name,
		) {
			state := states[i]
			return &state, nil
		}
	}

	return nil, errors.New("state not found")
}

// GetAll returns all states and the FCT.
func (r *StateJSONRepository) GetAll() []models.State {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	states, err := r.load()

	if err != nil {
		return []models.State{}
	}

	return states
}
