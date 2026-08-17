package repositories

import "nigeriaonlinevoting/models"

type VerificationRepository interface {
	Create(
		verification models.Verification,
	) error

	GetByUserID(
		userID int,
	) (*models.Verification, error)

	GetAll() []models.Verification

	Update(
		verification models.Verification,
	) error
}
