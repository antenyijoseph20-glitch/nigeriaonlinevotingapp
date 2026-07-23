package models

import "time"

type Party struct {
	ID int

	// Official party name
	Name string

	// Short abbreviation
	Abbreviation string

	// Party slogan
	Slogan string

	// Party logo filename or URL
	Logo string

	// Party chairman
	Chairman string

	// Party headquarters
	Headquarters string

	// Brief description
	Description string

	// Indicates whether the party
	// is approved to participate
	IsActive bool

	CreatedAt time.Time

	UpdatedAt time.Time
}
