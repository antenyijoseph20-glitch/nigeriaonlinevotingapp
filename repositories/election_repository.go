package repositories

import (
	"errors"
	"sync"

	"nigeriaonlinevoting/models"
)

type ElectionRepository interface {

	// Create a new election
	Create(
		election models.Election,
	) error

	// Update an election
	Update(
		election models.Election,
	) error

	// Get election by ID
	GetByID(
		id int,
	) (*models.Election, error)

	// Get the currently active election
	GetCurrent() (*models.Election, error)

	// Get all elections
	GetAll() []models.Election

	// Delete an election
	Delete(
		id int,
	) error
}

type ElectionJSONRepository struct {
	filePath string
	mutex    sync.Mutex
}

// Constructor
func NewElectionRepository(
	filePath string,
) *ElectionJSONRepository {

	return &ElectionJSONRepository{
		filePath: filePath,
	}
}

// Create

func (r *ElectionJSONRepository) Create(
	election models.Election,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	elections, err := r.load()

	if err != nil {
		return err
	}

	election.ID = len(elections) + 1

	elections = append(
		elections,
		election,
	)

	return r.save(
		elections,
	)
}

// Update

func (r *ElectionJSONRepository) Update(
	election models.Election,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	elections, err := r.load()

	if err != nil {
		return err
	}

	for i := range elections {

		if elections[i].ID == election.ID {

			elections[i] = election

			return r.save(
				elections,
			)
		}
	}

	return errors.New(
		"election not found",
	)
}

// GetByID

func (r *ElectionJSONRepository) GetByID(
	id int,
) (*models.Election, error) {

	elections, err := r.load()

	if err != nil {
		return nil, err
	}

	for _, election := range elections {

		if election.ID == id {

			return &election, nil
		}
	}

	return nil, errors.New(
		"election not found",
	)
}

// GetCurrent

func (r *ElectionJSONRepository) GetCurrent() (*models.Election, error) {

	elections, err := r.load()

	if err != nil {
		return nil, err
	}

	for _, election := range elections {

		if election.Status == "open" {

			return &election, nil
		}
	}

	return nil, errors.New(
		"no active election",
	)
}

// GetAll

func (r *ElectionJSONRepository) GetAll() []models.Election {

	elections, err := r.load()

	if err != nil {
		return []models.Election{}
	}

	return elections
}

// Delete

func (r *ElectionJSONRepository) Delete(
	id int,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	elections, err := r.load()

	if err != nil {
		return err
	}

	for i := range elections {

		if elections[i].ID == id {

			elections = append(
				elections[:i],
				elections[i+1:]...,
			)

			return r.save(
				elections,
			)
		}
	}

	return errors.New(
		"election not found",
	)
}
