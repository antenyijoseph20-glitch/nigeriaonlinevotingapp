package models

import "time"

type Election struct {
	ID int

	Title string

	Description string

	StartDate time.Time

	EndDate time.Time

	Status string

	CreatedAt time.Time

	UpdatedAt time.Time
}
