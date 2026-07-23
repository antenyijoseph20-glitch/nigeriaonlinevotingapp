package repositories

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"nigeriaonlinevoting/models"
)

type CandidateRepository interface {

	// Create
	Create(
		candidate models.Candidate,
	) error

	// Update
	Update(
		candidate models.Candidate,
	) error

	// Delete
	Delete(
		id int,
	) error

	// Get candidate by ID
	GetByID(
		id int,
	) (*models.Candidate, error)

	// Get all candidates
	GetAll() []models.Candidate

	// Get candidates by election
	GetByElectionID(
		electionID int,
	) []models.Candidate

	// Get candidates by party
	GetByPartyID(
		partyID int,
	) []models.Candidate

	// Get candidates by position
	GetByPositionID(
		positionID int,
	) []models.Candidate
}

type CandidateJSONRepository struct {
	filePath string
	mutex    sync.Mutex
}

// =====================================
// Constructor
// =====================================

func NewCandidateRepository(
	filePath string,
) *CandidateJSONRepository {

	return &CandidateJSONRepository{
		filePath: filePath,
	}
}

// =====================================
// Load
// =====================================

func (r *CandidateJSONRepository) load() ([]models.Candidate, error) {

	file, err := os.Open(r.filePath)

	if os.IsNotExist(err) {

		return []models.Candidate{}, nil
	}

	if err != nil {

		return nil, err
	}

	defer file.Close()

	var candidates []models.Candidate

	err = json.NewDecoder(file).Decode(
		&candidates,
	)

	if err != nil {

		return []models.Candidate{}, nil
	}

	return candidates, nil
}

// =====================================
// Save
// =====================================

func (r *CandidateJSONRepository) save(
	candidates []models.Candidate,
) error {

	file, err := os.Create(
		r.filePath,
	)

	if err != nil {

		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	encoder.SetIndent("", "    ")

	return encoder.Encode(
		candidates,
	)
}

// =====================================
// Create
// =====================================

func (r *CandidateJSONRepository) Create(
	candidate models.Candidate,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	candidates, err := r.load()

	if err != nil {

		return err
	}

	candidate.ID = len(candidates) + 1

	candidates = append(
		candidates,
		candidate,
	)

	return r.save(
		candidates,
	)
}

// =====================================
// Update
// =====================================

func (r *CandidateJSONRepository) Update(
	candidate models.Candidate,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	candidates, err := r.load()

	if err != nil {

		return err
	}

	for i := range candidates {

		if candidates[i].ID == candidate.ID {

			candidates[i] = candidate

			return r.save(
				candidates,
			)
		}
	}

	return errors.New(
		"candidate not found",
	)
}

// =====================================
// Delete
// =====================================

func (r *CandidateJSONRepository) Delete(
	id int,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	candidates, err := r.load()

	if err != nil {

		return err
	}

	for i := range candidates {

		if candidates[i].ID == id {

			candidates = append(
				candidates[:i],
				candidates[i+1:]...,
			)

			return r.save(
				candidates,
			)
		}
	}

	return errors.New(
		"candidate not found",
	)
}

// =====================================
// Get By ID
// =====================================

func (r *CandidateJSONRepository) GetByID(
	id int,
) (*models.Candidate, error) {

	candidates, err := r.load()

	if err != nil {

		return nil, err
	}

	for _, candidate := range candidates {

		if candidate.ID == id {

			return &candidate, nil
		}
	}

	return nil, errors.New(
		"candidate not found",
	)
}

// =====================================
// Get All
// =====================================

func (r *CandidateJSONRepository) GetAll() []models.Candidate {

	candidates, err := r.load()

	if err != nil {

		return []models.Candidate{}
	}

	return candidates
}

// =====================================
// Get By Election ID
// =====================================

func (r *CandidateJSONRepository) GetByElectionID(
	electionID int,
) []models.Candidate {

	candidates, _ := r.load()

	var result []models.Candidate

	for _, candidate := range candidates {

		if candidate.ElectionID == electionID {

			result = append(
				result,
				candidate,
			)
		}
	}

	return result
}

// =====================================
// Get By Party ID
// =====================================

func (r *CandidateJSONRepository) GetByPartyID(
	partyID int,
) []models.Candidate {

	candidates, _ := r.load()

	var result []models.Candidate

	for _, candidate := range candidates {

		if candidate.PartyID == partyID {

			result = append(
				result,
				candidate,
			)
		}
	}

	return result
}

// =====================================
// Get By Position ID
// =====================================

func (r *CandidateJSONRepository) GetByPositionID(
	positionID int,
) []models.Candidate {

	candidates, _ := r.load()

	var result []models.Candidate

	for _, candidate := range candidates {

		if candidate.PositionID == positionID {

			result = append(
				result,
				candidate,
			)
		}
	}

	return result
}
