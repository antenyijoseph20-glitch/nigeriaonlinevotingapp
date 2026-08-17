package repositories

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"

	"nigeriaonlinevoting/models"
)

// LGARepository defines operations for managing
// Local Government Areas and FCT Area Councils.
type LGARepository interface {
	Create(lga models.LGA) error
	Update(lga models.LGA) error
	Delete(id int) error

	GetByID(id int) (*models.LGA, error)
	GetByCode(code int) (*models.LGA, error)

	GetAll() []models.LGA
	GetByStateID(stateID int) []models.LGA
	GetByStateCode(stateCode string) []models.LGA
	GetByName(name string) (*models.LGA, error)
}

// LGAJSONRepository stores LGA records in a JSON file.
type LGAJSONRepository struct {
	filePath string
	mutex    sync.RWMutex
}

// NewLGARepository creates a new JSON-backed LGA repository.
func NewLGARepository(filePath string) *LGAJSONRepository {
	return &LGAJSONRepository{
		filePath: filePath,
	}
}

// load reads all LGA records from the JSON file.
func (r *LGAJSONRepository) load() ([]models.LGA, error) {
	file, err := os.Open(r.filePath)

	if os.IsNotExist(err) {
		return []models.LGA{}, nil
	}

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lgas []models.LGA

	if err := json.NewDecoder(file).Decode(&lgas); err != nil {
		return nil, err
	}

	if lgas == nil {
		return []models.LGA{}, nil
	}

	return lgas, nil
}

// save writes all LGA records to the JSON file.
func (r *LGAJSONRepository) save(lgas []models.LGA) error {
	file, err := os.Create(r.filePath)

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	return encoder.Encode(lgas)
}

// Create adds a new LGA or Area Council.
func (r *LGAJSONRepository) Create(lga models.LGA) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	lga.Name = strings.TrimSpace(lga.Name)
	lga.StateCode = strings.TrimSpace(lga.StateCode)
	lga.Type = strings.TrimSpace(lga.Type)

	if lga.Name == "" {
		return errors.New("LGA name is required")
	}

	if lga.StateID <= 0 {
		return errors.New("invalid state ID")
	}

	if lga.Code <= 0 {
		return errors.New("invalid LGA code")
	}

	if lga.StateCode == "" {
		return errors.New("state code is required")
	}

	if lga.Type != "LGA" && lga.Type != "AREA_COUNCIL" {
		return errors.New("invalid LGA type")
	}

	lgas, err := r.load()

	if err != nil {
		return err
	}

	for _, existing := range lgas {
		if existing.Code == lga.Code {
			return errors.New("LGA source code already exists")
		}
	}

	for _, existing := range lgas {
		if existing.StateID == lga.StateID &&
			strings.EqualFold(
				strings.TrimSpace(existing.Name),
				lga.Name,
			) {
			return errors.New("LGA name already exists in this state")
		}
	}

	maxID := 0

	for _, existing := range lgas {
		if existing.ID > maxID {
			maxID = existing.ID
		}
	}

	lga.ID = maxID + 1

	lgas = append(lgas, lga)

	return r.save(lgas)
}

// Update modifies an existing LGA.
func (r *LGAJSONRepository) Update(lga models.LGA) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	lga.Name = strings.TrimSpace(lga.Name)
	lga.StateCode = strings.TrimSpace(lga.StateCode)
	lga.Type = strings.TrimSpace(lga.Type)

	if lga.ID <= 0 {
		return errors.New("invalid LGA ID")
	}

	if lga.Name == "" {
		return errors.New("LGA name is required")
	}

	if lga.StateID <= 0 {
		return errors.New("invalid state ID")
	}

	if lga.Code <= 0 {
		return errors.New("invalid LGA code")
	}

	if lga.StateCode == "" {
		return errors.New("state code is required")
	}

	if lga.Type != "LGA" && lga.Type != "AREA_COUNCIL" {
		return errors.New("invalid LGA type")
	}

	lgas, err := r.load()

	if err != nil {
		return err
	}

	targetIndex := -1

	for i := range lgas {
		if lgas[i].ID == lga.ID {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		return errors.New("LGA not found")
	}

	for i, existing := range lgas {
		if i == targetIndex {
			continue
		}

		if existing.Code == lga.Code {
			return errors.New("LGA source code already exists")
		}

		if existing.StateID == lga.StateID &&
			strings.EqualFold(
				strings.TrimSpace(existing.Name),
				lga.Name,
			) {
			return errors.New("LGA name already exists in this state")
		}
	}

	lgas[targetIndex] = lga

	return r.save(lgas)
}

// Delete removes an LGA by internal ID.
func (r *LGAJSONRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if id <= 0 {
		return errors.New("invalid LGA ID")
	}

	lgas, err := r.load()

	if err != nil {
		return err
	}

	for i := range lgas {
		if lgas[i].ID == id {
			lgas = append(lgas[:i], lgas[i+1:]...)
			return r.save(lgas)
		}
	}

	return errors.New("LGA not found")
}

// GetByID returns an LGA by internal database ID.
func (r *LGAJSONRepository) GetByID(id int) (*models.LGA, error) {
	if id <= 0 {
		return nil, errors.New("invalid LGA ID")
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	lgas, err := r.load()

	if err != nil {
		return nil, err
	}

	for i := range lgas {
		if lgas[i].ID == id {
			lga := lgas[i]
			return &lga, nil
		}
	}

	return nil, errors.New("LGA not found")
}

// GetByCode returns an LGA by its source/reference code.
func (r *LGAJSONRepository) GetByCode(code int) (*models.LGA, error) {
	if code <= 0 {
		return nil, errors.New("invalid LGA code")
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	lgas, err := r.load()

	if err != nil {
		return nil, err
	}

	for i := range lgas {
		if lgas[i].Code == code {
			lga := lgas[i]
			return &lga, nil
		}
	}

	return nil, errors.New("LGA not found")
}

// GetAll returns every LGA and Area Council.
func (r *LGAJSONRepository) GetAll() []models.LGA {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	lgas, err := r.load()

	if err != nil {
		return []models.LGA{}
	}

	return lgas
}

// GetByStateID returns all LGAs belonging to a state or FCT.
func (r *LGAJSONRepository) GetByStateID(stateID int) []models.LGA {
	if stateID <= 0 {
		return []models.LGA{}
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	lgas, err := r.load()

	if err != nil {
		return []models.LGA{}
	}

	result := make([]models.LGA, 0)

	for _, lga := range lgas {
		if lga.StateID == stateID {
			result = append(result, lga)
		}
	}

	return result
}

// GetByStateCode returns all LGAs belonging to a state code.
func (r *LGAJSONRepository) GetByStateCode(
	stateCode string,
) []models.LGA {
	stateCode = strings.TrimSpace(stateCode)

	if stateCode == "" {
		return []models.LGA{}
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	lgas, err := r.load()

	if err != nil {
		return []models.LGA{}
	}

	result := make([]models.LGA, 0)

	for _, lga := range lgas {
		if strings.EqualFold(
			strings.TrimSpace(lga.StateCode),
			stateCode,
		) {
			result = append(result, lga)
		}
	}

	return result
}

// GetByName returns an LGA by name.
func (r *LGAJSONRepository) GetByName(
	name string,
) (*models.LGA, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errors.New("LGA name is required")
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	lgas, err := r.load()

	if err != nil {
		return nil, err
	}

	for i := range lgas {
		if strings.EqualFold(
			strings.TrimSpace(lgas[i].Name),
			name,
		) {
			lga := lgas[i]
			return &lga, nil
		}
	}

	return nil, errors.New("LGA not found")
}
