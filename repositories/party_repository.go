package repositories

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"nigeriaonlinevoting/models"
)

type PartyRepository interface {

	// Create a political party
	Create(
		party models.Party,
	) error

	// Update a political party
	Update(
		party models.Party,
	) error

	// Delete a political party
	Delete(
		id int,
	) error

	// Find by ID
	GetByID(
		id int,
	) (*models.Party, error)

	// Find by name
	GetByName(
		name string,
	) (*models.Party, error)

	// Return all parties
	GetAll() []models.Party
}

type PartyJSONRepository struct {
	filePath string
	mutex    sync.Mutex
}

// Constructor

func NewPartyRepository(
	filePath string,
) *PartyJSONRepository {

	return &PartyJSONRepository{
		filePath: filePath,
	}
}

//
// PRIVATE METHODS
//

func (r *PartyJSONRepository) load() ([]models.Party, error) {

	var parties []models.Party

	file, err := os.Open(r.filePath)

	if os.IsNotExist(err) {

		return parties, nil
	}

	if err != nil {

		return nil, err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(
		&parties,
	)

	if err != nil {

		return nil, err
	}

	return parties, nil
}

func (r *PartyJSONRepository) save(
	parties []models.Party,
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
		parties,
	)
}

//
// PUBLIC METHODS
//

func (r *PartyJSONRepository) Create(
	party models.Party,
) error {

	r.mutex.Lock()

	defer r.mutex.Unlock()

	parties, err := r.load()

	if err != nil {

		return err
	}

	party.ID = len(parties) + 1

	parties = append(
		parties,
		party,
	)

	return r.save(
		parties,
	)
}

func (r *PartyJSONRepository) Update(
	party models.Party,
) error {

	r.mutex.Lock()

	defer r.mutex.Unlock()

	parties, err := r.load()

	if err != nil {

		return err
	}

	for i := range parties {

		if parties[i].ID == party.ID {

			parties[i] = party

			return r.save(
				parties,
			)
		}
	}

	return errors.New(
		"party not found",
	)
}

func (r *PartyJSONRepository) Delete(
	id int,
) error {

	r.mutex.Lock()

	defer r.mutex.Unlock()

	parties, err := r.load()

	if err != nil {

		return err
	}

	for i := range parties {

		if parties[i].ID == id {

			parties = append(
				parties[:i],
				parties[i+1:]...,
			)

			return r.save(
				parties,
			)
		}
	}

	return errors.New(
		"party not found",
	)
}

func (r *PartyJSONRepository) GetByID(
	id int,
) (*models.Party, error) {

	parties, err := r.load()

	if err != nil {

		return nil, err
	}

	for _, party := range parties {

		if party.ID == id {

			return &party, nil
		}
	}

	return nil, errors.New(
		"party not found",
	)
}

func (r *PartyJSONRepository) GetByName(
	name string,
) (*models.Party, error) {

	parties, err := r.load()

	if err != nil {

		return nil, err
	}

	for _, party := range parties {

		if party.Name == name {

			return &party, nil
		}
	}

	return nil, errors.New(
		"party not found",
	)
}

func (r *PartyJSONRepository) GetAll() []models.Party {

	parties, err := r.load()

	if err != nil {

		return []models.Party{}
	}

	return parties
}
