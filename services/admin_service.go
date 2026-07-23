package services

import (
	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/repositories"
)

type AdminService struct {
	userRepo         repositories.UserRepository
	verificationRepo repositories.VerificationRepository
}

// NewAdminService creates a new AdminService.
func NewAdminService(
	userRepo repositories.UserRepository,
	verificationRepo repositories.VerificationRepository,
) *AdminService {

	return &AdminService{
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
	}
}

// GetDashboardStatistics returns statistics for the admin dashboard.
func (s *AdminService) GetDashboardStatistics() models.AdminStatistics {

	stats := models.AdminStatistics{}

	// ============================================
	// User Statistics
	// ============================================

	users := s.userRepo.GetAll()

	stats.TotalUsers = len(users)

	for _, user := range users {

		switch user.Role {

		case "voter":
			stats.TotalVoters++

		case "admin", "super_admin":
			stats.TotalAdmins++
		}

		if user.IsVerified {
			stats.VerifiedVoters++
		}

		if user.HasVoted {
			stats.TotalVotesCast++
		}
	}

	// ============================================
	// Verification Statistics
	// ============================================

	verifications := s.verificationRepo.GetAll()

	for _, verification := range verifications {

		switch verification.Status {

		case "pending":
			stats.PendingVerifications++

		case "approved":
			stats.ApprovedVerifications++

		case "rejected":
			stats.RejectedVerifications++
		}
	}

	return stats
}

// GetAllVerifications returns every verification request.
func (s *AdminService) GetAllVerifications() []models.Verification {
	return s.verificationRepo.GetAll()
}

// GetPendingVerifications returns only pending requests.
func (s *AdminService) GetPendingVerifications() []models.Verification {

	all := s.verificationRepo.GetAll()

	var pending []models.Verification

	for _, verification := range all {

		if verification.Status == "pending" {
			pending = append(pending, verification)
		}
	}

	return pending
}

// GetApprovedVerifications returns all approved verification requests.
func (s *AdminService) GetApprovedVerifications() []models.Verification {

	all := s.verificationRepo.GetAll()

	var approved []models.Verification

	for _, verification := range all {

		if verification.Status == "approved" {
			approved = append(approved, verification)
		}
	}

	return approved
}

// GetRejectedVerifications returns all rejected verification requests.
func (s *AdminService) GetRejectedVerifications() []models.Verification {

	all := s.verificationRepo.GetAll()

	var rejected []models.Verification

	for _, verification := range all {

		if verification.Status == "rejected" {
			rejected = append(rejected, verification)
		}
	}

	return rejected
}
