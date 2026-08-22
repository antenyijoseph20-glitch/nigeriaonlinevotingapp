package repositories

import "nigeriaonlinevoting/models"

// ElectoralVerificationRepository defines storage operations
// for electoral verification records.
type ElectoralVerificationRepository interface {
	Create(
		verification models.ElectoralVerification,
	) error

	GetByUserID(
		userID int,
	) (*models.ElectoralVerification, error)

	GetByID(
		id int,
	) (*models.ElectoralVerification, error)

	GetAll() []models.ElectoralVerification

	Update(
		verification models.ElectoralVerification,
	) error
}