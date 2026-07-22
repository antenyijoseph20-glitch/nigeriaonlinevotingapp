package models

import "time"


type Verification struct {

	ID int

	UserID int


	FullName string

	DateOfBirth string

	Gender string


	State string

	LGA string

	Ward string

	PollingUnit string


	Status string
	// pending
	// approved
	// rejected


	SubmittedAt time.Time

	ReviewedAt time.Time

	ReviewedBy int
}