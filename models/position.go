package models

import "time"

// Position represents an elective office
// that candidates can contest for.
type Position struct {
	// Unique identifier.
	ID int `json:"id"`

	// Name of the position.
	// Example:
	// President
	// Governor
	// Senator
	// Councillor
	Name string `json:"name"`

	// Brief description of the position.
	Description string `json:"description"`

	// Government level.
	// Federal
	// State
	// Local
	Level string `json:"level"`

	// Number of available seats.
	// Example:
	// President = 1
	// Governor = 1
	// Senator = 3
	Seats int `json:"seats"`

	// Whether the position is active.
	IsActive bool `json:"is_active"`

	// Record timestamps.
	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}
