package repositories

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"nigeriaonlinevoting/models"
)

// VoterRepository defines operations
// for managing registered voters.
type VoterRepository interface {

	// Create a voter.
	Create(
		voter models.Voter,
	) error

	// Update a voter.
	Update(
		voter models.Voter,
	) error

	// Delete a voter.
	Delete(
		id int,
	) error

	// Get voter by ID.
	GetByID(
		id int,
	) (*models.Voter, error)

	// Get voter by linked user ID.
	GetByUserID(
		userID int,
	) (*models.Voter, error)

	// Get voter by NIN.
	GetByNIN(
		nin string,
	) (*models.Voter, error)

	// Get voter by VIN.
	GetByVIN(
		vin string,
	) (*models.Voter, error)

	// Get voter by PVC number.
	GetByPVCNumber(
		pvcNumber string,
	) (*models.Voter, error)

	// Return all voters.
	GetAll() []models.Voter
}

// =====================================
// JSON Repository
// =====================================

type VoterJSONRepository struct {
	filePath string
	mutex    sync.Mutex
}

// =====================================
// Constructor
// =====================================

func NewVoterRepository(
	filePath string,
) *VoterJSONRepository {

	return &VoterJSONRepository{
		filePath: filePath,
	}
}

// =====================================
// Load
// =====================================

func (r *VoterJSONRepository) load() ([]models.Voter, error) {

	file, err := os.Open(r.filePath)

	if os.IsNotExist(err) {
		return []models.Voter{}, nil
	}

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var voters []models.Voter

	err = json.NewDecoder(file).Decode(&voters)

	if err != nil {
		return nil, err
	}

	return voters, nil
}

// =====================================
// Save
// =====================================

func (r *VoterJSONRepository) save(
	voters []models.Voter,
) error {

	file, err := os.Create(r.filePath)

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	encoder.SetIndent("", "    ")

	return encoder.Encode(voters)
}

// =====================================
// Create
// =====================================

func (r *VoterJSONRepository) Create(
	voter models.Voter,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if voter.UserID <= 0 {
		return errors.New("invalid user ID")
	}

	if voter.NIN == "" {
		return errors.New("NIN is required")
	}

	if voter.VIN == "" {
		return errors.New("VIN is required")
	}

	if voter.PVCNumber == "" {
		return errors.New("PVC number is required")
	}

	voters, err := r.load()

	if err != nil {
		return err
	}

	// Prevent duplicate UserID.
	for _, existing := range voters {

		if existing.UserID == voter.UserID {
			return errors.New(
				"voter already exists for this user",
			)
		}
	}

	// Prevent duplicate NIN.
	for _, existing := range voters {

		if existing.NIN == voter.NIN {
			return errors.New(
				"NIN already exists",
			)
		}
	}

	// Prevent duplicate VIN.
	for _, existing := range voters {

		if existing.VIN == voter.VIN {
			return errors.New(
				"VIN already exists",
			)
		}
	}

	// Prevent duplicate PVC number.
	for _, existing := range voters {

		if existing.PVCNumber == voter.PVCNumber {
			return errors.New(
				"PVC number already exists",
			)
		}
	}

	// Generate the next safe ID.
	maxID := 0

	for _, existing := range voters {

		if existing.ID > maxID {
			maxID = existing.ID
		}
	}

	voter.ID = maxID + 1

	voters = append(
		voters,
		voter,
	)

	return r.save(voters)
}

// =====================================
// Update
// =====================================

func (r *VoterJSONRepository) Update(
	voter models.Voter,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if voter.ID <= 0 {
		return errors.New("invalid voter ID")
	}

	if voter.UserID <= 0 {
		return errors.New("invalid user ID")
	}

	if voter.NIN == "" {
		return errors.New("NIN is required")
	}

	if voter.VIN == "" {
		return errors.New("VIN is required")
	}

	if voter.PVCNumber == "" {
		return errors.New("PVC number is required")
	}

	voters, err := r.load()

	if err != nil {
		return err
	}

	targetIndex := -1

	for i := range voters {

		if voters[i].ID == voter.ID {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		return errors.New("voter not found")
	}

	// Check for duplicate UserID.
	for i, existing := range voters {

		if i != targetIndex &&
			existing.UserID == voter.UserID {

			return errors.New(
				"user already linked to another voter",
			)
		}
	}

	// Check for duplicate NIN.
	for i, existing := range voters {

		if i != targetIndex &&
			existing.NIN == voter.NIN {

			return errors.New(
				"NIN already belongs to another voter",
			)
		}
	}

	// Check for duplicate VIN.
	for i, existing := range voters {

		if i != targetIndex &&
			existing.VIN == voter.VIN {

			return errors.New(
				"VIN already belongs to another voter",
			)
		}
	}

	// Check for duplicate PVC number.
	for i, existing := range voters {

		if i != targetIndex &&
			existing.PVCNumber == voter.PVCNumber {

			return errors.New(
				"PVC number already belongs to another voter",
			)
		}
	}

	voters[targetIndex] = voter

	return r.save(voters)
}

// =====================================
// Delete
// =====================================

func (r *VoterJSONRepository) Delete(
	id int,
) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if id <= 0 {
		return errors.New("invalid voter ID")
	}

	voters, err := r.load()

	if err != nil {
		return err
	}

	for i := range voters {

		if voters[i].ID == id {

			voters = append(
				voters[:i],
				voters[i+1:]...,
			)

			return r.save(voters)
		}
	}

	return errors.New("voter not found")
}

// =====================================
// Get By ID
// =====================================

func (r *VoterJSONRepository) GetByID(
	id int,
) (*models.Voter, error) {

	if id <= 0 {
		return nil, errors.New(
			"invalid voter ID",
		)
	}

	voters, err := r.load()

	if err != nil {
		return nil, err
	}

	for i := range voters {

		if voters[i].ID == id {
			return &voters[i], nil
		}
	}

	return nil, errors.New(
		"voter not found",
	)
}

// =====================================
// Get By User ID
// =====================================

func (r *VoterJSONRepository) GetByUserID(
	userID int,
) (*models.Voter, error) {

	if userID <= 0 {
		return nil, errors.New(
			"invalid user ID",
		)
	}

	voters, err := r.load()

	if err != nil {
		return nil, err
	}

	for i := range voters {

		if voters[i].UserID == userID {
			return &voters[i], nil
		}
	}

	return nil, errors.New(
		"voter not found",
	)
}

// =====================================
// Get By NIN
// =====================================

func (r *VoterJSONRepository) GetByNIN(
	nin string,
) (*models.Voter, error) {

	if nin == "" {
		return nil, errors.New(
			"NIN is required",
		)
	}

	voters, err := r.load()

	if err != nil {
		return nil, err
	}

	for i := range voters {

		if voters[i].NIN == nin {
			return &voters[i], nil
		}
	}

	return nil, errors.New(
		"voter not found",
	)
}

// =====================================
// Get By VIN
// =====================================

func (r *VoterJSONRepository) GetByVIN(
	vin string,
) (*models.Voter, error) {

	if vin == "" {
		return nil, errors.New(
			"VIN is required",
		)
	}

	voters, err := r.load()

	if err != nil {
		return nil, err
	}

	for i := range voters {

		if voters[i].VIN == vin {
			return &voters[i], nil
		}
	}

	return nil, errors.New(
		"voter not found",
	)
}

// =====================================
// Get By PVC Number
// =====================================

func (r *VoterJSONRepository) GetByPVCNumber(
	pvcNumber string,
) (*models.Voter, error) {

	if pvcNumber == "" {
		return nil, errors.New(
			"PVC number is required",
		)
	}

	voters, err := r.load()

	if err != nil {
		return nil, err
	}

	for i := range voters {

		if voters[i].PVCNumber == pvcNumber {
			return &voters[i], nil
		}
	}

	return nil, errors.New(
		"voter not found",
	)
}

// =====================================
// Get All
// =====================================

func (r *VoterJSONRepository) GetAll() []models.Voter {

	voters, err := r.load()

	if err != nil {
		return []models.Voter{}
	}

	return voters
}
