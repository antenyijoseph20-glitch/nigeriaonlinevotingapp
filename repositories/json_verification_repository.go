package repositories

import (
	"encoding/json"
	"errors"
	"os"

	"nigeriaonlinevoting/models"
)

type JSONVerificationRepository struct {
	filePath string
}

func NewVerificationRepository(
	path string,
) *JSONVerificationRepository {

	return &JSONVerificationRepository{
		filePath: path,
	}

}

func (j *JSONVerificationRepository) load() ([]models.Verification, error) {

	data, err := os.ReadFile(j.filePath)

	if err != nil {
		return nil, err
	}

	var verifications []models.Verification

	err = json.Unmarshal(
		data,
		&verifications,
	)

	return verifications, err

}

func (j *JSONVerificationRepository) save(
	verifications []models.Verification,
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

func (j *JSONVerificationRepository) Create(
	v models.Verification,
) error {

	verifications, err := j.load()

	if err != nil {
		return err
	}

	v.ID = len(verifications) + 1

	verifications = append(
		verifications,
		v,
	)

	return j.save(verifications)

}

func (j *JSONVerificationRepository) GetByUserID(
	userID int,
) (*models.Verification, error) {

	verifications, err := j.load()

	if err != nil {
		return nil, err
	}

	for _, v := range verifications {

		if v.UserID == userID {

			return &v, nil

		}
	}

	return nil, errors.New("verification not found")

}

func (j *JSONVerificationRepository) GetAll() []models.Verification {

	verifications, err := j.load()

	if err != nil {
		return []models.Verification{}
	}

	return verifications

}

func (j *JSONVerificationRepository) Update(
	updated models.Verification,
) error {

	verifications, err := j.load()

	if err != nil {
		return err
	}

	for i, v := range verifications {

		if v.ID == updated.ID {

			verifications[i] = updated

			return j.save(verifications)

		}

	}

	return errors.New("verification not found")

}
