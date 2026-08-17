package models

import "time"

// LGA represents a Local Government Area in one of
// Nigeria's 36 states or an Area Council in the FCT.
type LGA struct {
	ID int

	// Source reference number used by the
	// verified dataset.
	// This is separate from our internal ID.
	Code int

	// ID of the parent State record.
	StateID int

	// State abbreviation.
	// Examples:
	// AB = Abia
	// AD = Adamawa
	// FC = Federal Capital Territory
	StateCode string

	// Verified LGA or FCT Area Council name.
	Name string

	// Administrative type.
	//
	// LGA = Local Government Area
	// AREA_COUNCIL = FCT Area Council
	Type string

	// Whether this administrative record is active.
	IsActive bool

	// Record timestamps.
	CreatedAt time.Time
	UpdatedAt time.Time
}
