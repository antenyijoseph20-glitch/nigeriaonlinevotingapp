package repositories

import (
	"encoding/json"
	"errors"
	"os"

	"nigeriaonlinevoting/models"
)

// JSONElectoralVerificationRepository stores electoral
// verification records in a JSON file.
type JSONElectoralVerificationRepository struct {
	filePath string
}

// NewElectoralVerificationRepository creates a new JSON
// electoral verification repository.
func NewElectoralVerificationRepository(
	path string,
) *JSONElectoralVerificationRepository {

	return &JSONElectoralVerificationRepository{
		filePath: path,
	}
}

// load reads all electoral verification records from disk.
func (j *JSONElectoralVerificationRepository) load() (
	[]models.ElectoralVerification,
	error,
) {

	data, err := os.ReadFile(j.filePath)

	if err != nil {
		return nil, err
	}

	var verifications []models.ElectoralVerification

	err = json.Unmarshal(
		data,
		&verifications,
	)

	if err != nil {
		return nil, err
	}

	return verifications, nil
}

// save writes all electoral verification records to disk.
func (j *JSONElectoralVerificationRepository) save(
	verifications []models.ElectoralVerification,
) error {

	data, err := json.MarshalIndent(
		verifications,
		"",
		"    ",
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		j.filePath,
		data,
		0644,
	)
}

// Create adds a new electoral verification record.
func (j *JSONElectoralVerificationRepository) Create(
	verification models.ElectoralVerification,
) error {

	verifications, err := j.load()

	if err != nil {
		return err
	}

	// Generate an ID greater than every existing ID.
	nextID := 1

	for _, existing := range verifications {

		if existing.ID >= nextID {
			nextID = existing.ID + 1
		}
	}

	verification.ID = nextID

	verifications = append(
		verifications,
		verification,
	)

	return j.save(verifications)
}

// GetByUserID returns the electoral verification record
// belonging to a specific user.
func (j *JSONElectoralVerificationRepository) GetByUserID(
	userID int,
) (*models.ElectoralVerification, error) {

	verifications, err := j.load()

	if err != nil {
		return nil, err
	}

	for _, verification := range verifications {

		if verification.UserID == userID {
			result := verification
			return &result, nil
		}
	}

	return nil, errors.New(
		"electoral verification not found",
	)
}

// GetByID returns an electoral verification record
// by its primary ID.
func (j *JSONElectoralVerificationRepository) GetByID(
	id int,
) (*models.ElectoralVerification, error) {

	verifications, err := j.load()

	if err != nil {
		return nil, err
	}

	for _, verification := range verifications {

		if verification.ID == id {
			result := verification
			return &result, nil
		}
	}

	return nil, errors.New(
		"electoral verification not found",
	)
}

// GetAll returns all electoral verification records.
func (j *JSONElectoralVerificationRepository) GetAll() []models.ElectoralVerification {

	verifications, err := j.load()

	if err != nil {
		return []models.ElectoralVerification{}
	}

	return verifications
}

// Update replaces an existing electoral verification record.
func (j *JSONElectoralVerificationRepository) Update(
	updated models.ElectoralVerification,
) error {

	verifications, err := j.load()

	if err != nil {
		return err
	}

	for i, verification := range verifications {

		if verification.ID == updated.ID {

			verifications[i] = updated

			return j.save(verifications)
		}
	}

	return errors.New(
		"electoral verification not found",
	)
}