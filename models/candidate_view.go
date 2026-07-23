package models

type CandidateView struct {
	ID int

	FirstName string

	LastName string

	Gender string

	Email string

	PhoneNumber string

	IsApproved bool

	IsActive bool

	ElectionName string

	PartyName string

	PositionName string
}
