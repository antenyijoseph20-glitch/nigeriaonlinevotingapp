package models

import "time"

// Position represents an elective office
// that candidates can contest for.
type Position struct {

	// Unique identifier
	ID int

	// Name of the position
	// Example:
	// President
	// Governor
	// Senator
	// Councillor
	Name string

	// Brief description of the position
	Description string

	// Government level
	// Federal
	// State
	// Local Government
	Level string

	// Number of available seats
	// Example:
	// President = 1
	// Governor = 1
	// Senator = 3
	Seats int

	// Whether the position is active
	IsActive bool

	// Record timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
}
