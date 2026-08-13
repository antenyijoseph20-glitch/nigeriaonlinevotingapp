package repositories

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	"nigeriaonlinevoting/models"
)

// =====================================
// Election Repository Interface
// =====================================

type ElectionRepository interface {

	// Create a new election.
	Create(
		election models.Election,
	) error

	// Update an election.
	Update(
		election models.Election,
	) error

	// Get election by ID.
	GetByID(
		id int,
	) (*models.Election, error)

	// Get the currently active election.
	GetCurrent() (*models.Election, error)

	// Get all elections.
	GetAll() []models.Election

	// Delete an election.
	Delete(
		id int,
	) error
}

// =====================================
// JSON Repository
// =====================================

type ElectionJSONRepository struct {
	filePath string
	mutex    sync.Mutex
}

// =====================================
// Constructor
// =====================================

func NewElectionRepository(
	filePath string,
) *ElectionJSONRepository {

	return &ElectionJSONRepository{
		filePath: filePath,
	}
}

// =====================================
// Load
// =====================================

func (r *ElectionJSONRepository) load() (
	[]models.Election,
	error,
) {

	file, err := os.Open(
		r.filePath,
	)

	if os.IsNotExist(err) {
		return []models.Election{}, nil
	}

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var elections []models.Election

	err = json.NewDecoder(file).Decode(
		&elections,
	)

	if err != nil {

		// An empty file should behave like
		// an empty repository.
		if errors.Is(
			err,
			io.EOF,
		) {
			return []models.Election{}, nil
		}

		return nil, err
	}

	if elections == nil {
		return []models.Election{}, nil
	}

	return elections, nil
}

// =====================================
// Save
// =====================================

func (r *ElectionJSONRepository) save(
	elections []models.Election,
) error {

	// Ensure the directory exists when
	// the repository receives a path
	// inside an existing directory.
	//
	// We deliberately do not create arbitrary
	// parent directories here because a bad
	// configuration should fail visibly.

	file, err := os.Create(
		r.filePath,
	)

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	encoder.SetIndent(
		"",
		"    ",
	)

	return encoder.Encode(
		elections,
	)
}

// =====================================
// Create
// =====================================

func (r *ElectionJSONRepository) Create(
	election models.Election,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// -------------------------------------
	// Validate basic data
	// -------------------------------------

	if election.Title == "" {
		return errors.New(
			"election title is required",
		)
	}

	if election.StartDate.IsZero() {
		return errors.New(
			"election start date is required",
		)
	}

	if election.EndDate.IsZero() {
		return errors.New(
			"election end date is required",
		)
	}

	if !election.EndDate.After(
		election.StartDate,
	) {
		return errors.New(
			"end date must be after start date",
		)
	}

	// -------------------------------------
	// Load existing elections
	// -------------------------------------

	elections, err := r.load()

	if err != nil {
		return err
	}

	// -------------------------------------
	// Prevent more than one open election.
	// -------------------------------------

	if election.Status == "open" {

		for _, existing := range elections {

			if existing.Status == "open" {

				return errors.New(
					"another election is already open",
				)
			}
		}
	}

	// -------------------------------------
	// Generate safe next ID.
	// -------------------------------------

	maxID := 0

	for _, existing := range elections {

		if existing.ID > maxID {
			maxID = existing.ID
		}
	}

	election.ID = maxID + 1

	elections = append(
		elections,
		election,
	)

	return r.save(
		elections,
	)
}

// =====================================
// Update
// =====================================

func (r *ElectionJSONRepository) Update(
	election models.Election,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if election.ID <= 0 {
		return errors.New(
			"invalid election ID",
		)
	}

	if election.Title == "" {
		return errors.New(
			"election title is required",
		)
	}

	if election.StartDate.IsZero() {
		return errors.New(
			"election start date is required",
		)
	}

	if election.EndDate.IsZero() {
		return errors.New(
			"election end date is required",
		)
	}

	if !election.EndDate.After(
		election.StartDate,
	) {
		return errors.New(
			"end date must be after start date",
		)
	}

	elections, err := r.load()

	if err != nil {
		return err
	}

	targetIndex := -1

	for i := range elections {

		if elections[i].ID == election.ID {

			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		return errors.New(
			"election not found",
		)
	}

	// -------------------------------------
	// Prevent multiple open elections.
	// -------------------------------------

	if election.Status == "open" {

		for i, existing := range elections {

			if i == targetIndex {
				continue
			}

			if existing.Status == "open" {

				return errors.New(
					"another election is already open",
				)
			}
		}
	}

	// -------------------------------------
	// Preserve immutable creation data.
	// -------------------------------------

	election.CreatedAt =
		elections[targetIndex].CreatedAt

	elections[targetIndex] = election

	return r.save(
		elections,
	)
}

// =====================================
// Get By ID
// =====================================

func (r *ElectionJSONRepository) GetByID(
	id int,
) (*models.Election, error) {

	if id <= 0 {
		return nil, errors.New(
			"invalid election ID",
		)
	}

	elections, err := r.load()

	if err != nil {
		return nil, err
	}

	for i := range elections {

		if elections[i].ID == id {

			return &elections[i], nil
		}
	}

	return nil, errors.New(
		"election not found",
	)
}

// =====================================
// Get Current
// =====================================

func (r *ElectionJSONRepository) GetCurrent() (
	*models.Election,
	error,
) {

	elections, err := r.load()

	if err != nil {
		return nil, err
	}

	for i := range elections {

		if elections[i].Status == "open" {

			return &elections[i], nil
		}
	}

	return nil, errors.New(
		"no active election",
	)
}

// =====================================
// Get All
// =====================================

func (r *ElectionJSONRepository) GetAll() []models.Election {

	elections, err := r.load()

	if err != nil {
		return []models.Election{}
	}

	return elections
}

// =====================================
// Delete
// =====================================

func (r *ElectionJSONRepository) Delete(
	id int,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if id <= 0 {
		return errors.New(
			"invalid election ID",
		)
	}

	elections, err := r.load()

	if err != nil {
		return err
	}

	for i := range elections {

		if elections[i].ID == id {

			if elections[i].Status == "open" {

				return errors.New(
					"cannot delete an open election",
				)
			}

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
