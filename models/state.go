package models

import "time"

// State represents one of Nigeria's 36 states
// or the Federal Capital Territory (FCT).
type State struct {
	ID int

	// Official name of the state or FCT.
	Name string

	// Short identifier used internally.
	// Examples:
	// AB = Abia
	// AD = Adamawa
	// FC = FCT
	Code string

	// Indicates whether this record is the
	// Federal Capital Territory.
	IsFCT bool

	// Whether the state record is currently active.
	IsActive bool

	// Record timestamps.
	CreatedAt time.Time
	UpdatedAt time.Time
}