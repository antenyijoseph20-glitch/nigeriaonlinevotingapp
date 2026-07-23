package models

import "time"

// Candidate represents a person contesting
// in an election.
type Candidate struct {

	// Unique candidate ID
	ID int

	// Election the candidate belongs to
	ElectionID int

	// Political party
	PartyID int

	// Position being contested
	PositionID int

	// Candidate details
	FirstName string
	LastName  string
	Gender    string

	DateOfBirth string

	Email string

	PhoneNumber string

	// Campaign information
	Biography string

	Manifesto string

	Photo string

	// Candidate status
	IsApproved bool

	IsActive bool

	CreatedAt time.Time

	UpdatedAt time.Time
}
