package repositories

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"nigeriaonlinevoting/models"
)

type JSONRepository struct {
	filePath string
}

func NewJSONRepository(path string) *JSONRepository {
	return &JSONRepository{
		filePath: path,
	}
}

// loadUsers reads all users from the JSON file.
func (j *JSONRepository) loadUsers() ([]models.User, error) {

	data, err := os.ReadFile(j.filePath)
	if err != nil {
		return nil, err
	}

	var users []models.User

	if len(data) == 0 {
		return users, nil
	}

	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// saveUsers writes all users back to the JSON file.
func (j *JSONRepository) saveUsers(users []models.User) error {

	data, err := json.MarshalIndent(users, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(j.filePath, data, 0644)
}

func (j *JSONRepository) Create(user models.User) error {

	users, err := j.loadUsers()
	if err != nil {
		return err
	}

	// Generate a new ID
	user.ID = len(users) + 1

	users = append(users, user)

	return j.saveUsers(users)
}

func (j *JSONRepository) GetByEmail(email string) (*models.User, error) {

	users, err := j.loadUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {

		if strings.EqualFold(user.Email, email) {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (j *JSONRepository) GetByID(id int) (*models.User, error) {

	users, err := j.loadUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {

		if user.ID == id {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (j *JSONRepository) Update(updatedUser models.User) error {

	users, err := j.loadUsers()
	if err != nil {
		return err
	}

	for i, user := range users {

		if user.ID == updatedUser.ID {
			users[i] = updatedUser
			return j.saveUsers(users)
		}
	}

	return errors.New("user not found")
}

func (j *JSONRepository) Delete(id int) error {

	users, err := j.loadUsers()
	if err != nil {
		return err
	}

	for i, user := range users {

		if user.ID == id {

			users = append(users[:i], users[i+1:]...)

			return j.saveUsers(users)
		}
	}

	return errors.New("user not found")
}

func (j *JSONRepository) GetAll() []models.User {

	users, err := j.loadUsers()
	if err != nil {
		return []models.User{}
	}

	return users
}
