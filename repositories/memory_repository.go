package repositories

import (
	"errors"
	"nigeriaonlinevoting/models"
)

type MemoryRepository struct {
	users  []models.User
	nextID int
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users:  []models.User{},
		nextID: 1,
	}
}

func (m *MemoryRepository) Create(user models.User) error {

	for _, u := range m.users {
		if u.Email == user.Email {
			return errors.New("email already exists")
		}
	}

	user.ID = m.nextID
	m.nextID++

	m.users = append(m.users, user)

	return nil
}

func (m *MemoryRepository) GetByEmail(email string) (*models.User, error) {

	for _, user := range m.users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (m *MemoryRepository) GetByID(id int) (*models.User, error) {

	for _, user := range m.users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (m *MemoryRepository) Update(user models.User) error {

	for i := range m.users {
		if m.users[i].ID == user.ID {
			m.users[i] = user
			return nil
		}
	}

	return errors.New("user not found")
}

func (m *MemoryRepository) Delete(id int) error {

	for i := range m.users {
		if m.users[i].ID == id {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}

	return errors.New("user not found")
}

func (m *MemoryRepository) GetAll() []models.User {
	return m.users
}
