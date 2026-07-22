package repositories

import "nigeriaonlinevoting/models"

// UserRepository defines the behavior required for user storage.
type UserRepository interface {
	Create(user models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id int) (*models.User, error)
	Update(user models.User) error
	Delete(id int) error
	GetAll() []models.User
}
