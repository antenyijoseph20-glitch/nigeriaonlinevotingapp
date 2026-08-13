package services

import (
	"errors"
	"strings"
	"time"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/repositories"
)

type ElectionService struct {
	ElectionRepo repositories.ElectionRepository
}

// =====================================
// Constructor
// =====================================

func NewElectionService(
	repo repositories.ElectionRepository,
) *ElectionService {

	return &ElectionService{
		ElectionRepo: repo,
	}
}

// =====================================
// Create Election
// =====================================

func (s *ElectionService) CreateElection(
	election models.Election,
) error {

	if s.ElectionRepo == nil {
		return errors.New("election repository is required")
	}

	// -------------------------------------
	// Clean input
	// -------------------------------------

	election.Title = strings.TrimSpace(
		election.Title,
	)

	election.Description = strings.TrimSpace(
		election.Description,
	)

	// -------------------------------------
	// Validate title
	// -------------------------------------

	if election.Title == "" {
		return errors.New(
			"election title is required",
		)
	}

	// -------------------------------------
	// Validate dates
	// -------------------------------------

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
	// Set initial status
	// -------------------------------------

	election.Status = "draft"

	now := time.Now()

	election.CreatedAt = now
	election.UpdatedAt = now

	return s.ElectionRepo.Create(
		election,
	)
}

// =====================================
// Get Election By ID
// =====================================

func (s *ElectionService) GetElectionByID(
	id int,
) (*models.Election, error) {

	if s.ElectionRepo == nil {
		return nil, errors.New(
			"election repository is required",
		)
	}

	if id <= 0 {
		return nil, errors.New(
			"invalid election ID",
		)
	}

	return s.ElectionRepo.GetByID(
		id,
	)
}

// =====================================
// Get Current Election
// =====================================

func (s *ElectionService) GetCurrentElection() (
	*models.Election,
	error,
) {

	if s.ElectionRepo == nil {
		return nil, errors.New(
			"election repository is required",
		)
	}

	return s.ElectionRepo.GetCurrent()
}

// =====================================
// Get All Elections
// =====================================

func (s *ElectionService) GetAllElections() []models.Election {

	if s.ElectionRepo == nil {
		return []models.Election{}
	}

	return s.ElectionRepo.GetAll()
}

// =====================================
// Update Election
// =====================================

func (s *ElectionService) UpdateElection(
	election models.Election,
) error {

	if s.ElectionRepo == nil {
		return errors.New(
			"election repository is required",
		)
	}

	if election.ID <= 0 {
		return errors.New(
			"invalid election ID",
		)
	}

	election.Title = strings.TrimSpace(
		election.Title,
	)

	election.Description = strings.TrimSpace(
		election.Description,
	)

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
	// Make sure election exists
	// -------------------------------------

	existing, err := s.ElectionRepo.GetByID(
		election.ID,
	)

	if err != nil {
		return err
	}

	// -------------------------------------
	// Do not allow an open election to
	// silently change its status.
	// -------------------------------------

	if existing.Status == "open" &&
		election.Status != "open" {

		return errors.New(
			"cannot change the status of an open election",
		)
	}

	election.CreatedAt = existing.CreatedAt

	election.UpdatedAt = time.Now()

	return s.ElectionRepo.Update(
		election,
	)
}

// =====================================
// Delete Election
// =====================================

func (s *ElectionService) DeleteElection(
	id int,
) error {

	if s.ElectionRepo == nil {
		return errors.New(
			"election repository is required",
		)
	}

	if id <= 0 {
		return errors.New(
			"invalid election ID",
		)
	}

	election, err := s.ElectionRepo.GetByID(
		id,
	)

	if err != nil {
		return err
	}

	if election.Status == "open" {
		return errors.New(
			"cannot delete an open election",
		)
	}

	return s.ElectionRepo.Delete(
		id,
	)
}

// =====================================
// Open Election
// =====================================

func (s *ElectionService) OpenElection(
	id int,
) error {

	if s.ElectionRepo == nil {
		return errors.New(
			"election repository is required",
		)
	}

	if id <= 0 {
		return errors.New(
			"invalid election ID",
		)
	}

	// -------------------------------------
	// Get target election first.
	// -------------------------------------

	election, err := s.ElectionRepo.GetByID(
		id,
	)

	if err != nil {
		return err
	}

	// -------------------------------------
	// Validate current status.
	// -------------------------------------

	if election.Status == "open" {
		return errors.New(
			"election is already open",
		)
	}

	if election.Status == "closed" {
		return errors.New(
			"closed election cannot be reopened",
		)
	}

	// -------------------------------------
	// Validate election dates.
	// -------------------------------------

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
			"election dates are invalid",
		)
	}

	now := time.Now()

	if now.Before(
		election.StartDate,
	) {
		return errors.New(
			"election cannot be opened before its start date",
		)
	}

	if !now.Before(
		election.EndDate,
	) {
		return errors.New(
			"election has already ended",
		)
	}

	// -------------------------------------
	// Only one election can be open.
	// -------------------------------------

	current, err := s.ElectionRepo.GetCurrent()

	if err != nil {
		return err
	}

	if current != nil &&
		current.ID != election.ID {

		return errors.New(
			"another election is already open",
		)
	}

	// -------------------------------------
	// Open election.
	// -------------------------------------

	election.Status = "open"
	election.UpdatedAt = now

	return s.ElectionRepo.Update(
		*election,
	)
}

// =====================================
// Close Election
// =====================================

func (s *ElectionService) CloseElection(
	id int,
) error {

	if s.ElectionRepo == nil {
		return errors.New(
			"election repository is required",
		)
	}

	if id <= 0 {
		return errors.New(
			"invalid election ID",
		)
	}

	election, err := s.ElectionRepo.GetByID(
		id,
	)

	if err != nil {
		return err
	}

	// -------------------------------------
	// Validate status.
	// -------------------------------------

	if election.Status == "closed" {
		return errors.New(
			"election is already closed",
		)
	}

	if election.Status != "open" {
		return errors.New(
			"only an open election can be closed",
		)
	}

	// -------------------------------------
	// Close election.
	// -------------------------------------

	election.Status = "closed"
	election.UpdatedAt = time.Now()

	return s.ElectionRepo.Update(
		*election,
	)
}
