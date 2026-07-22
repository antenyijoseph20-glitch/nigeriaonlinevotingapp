package models

import "time"

// User represents a registered user in the Nigeria Online Voting System.
type User struct {
	// ==========================
	// Primary Key
	// ==========================
	ID int `json:"id"`

	// ==========================
	// Personal Information
	// ==========================
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	OtherName string `json:"other_name"`

	DateOfBirth time.Time `json:"date_of_birth"`
	Gender      string    `json:"gender"`

	// ==========================
	// Contact Information
	// ==========================
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`

	// ==========================
	// Authentication
	// ==========================
	PasswordHash string `json:"password_hash"`

	// ==========================
	// National Identity
	// ==========================
	NIN string `json:"nin"`
	VIN string `json:"vin"`

	// ==========================
	// Voter Information
	// ==========================
	State       string `json:"state"`
	LGA         string `json:"lga"`
	Ward        string `json:"ward"`
	PollingUnit string `json:"polling_unit"`

	// ==========================
	// Account Status
	// ==========================
	Role          string `json:"role"`
	IsVerified    bool   `json:"is_verified"`
	HasVoted      bool   `json:"has_voted"`
	AccountActive bool   `json:"account_active"`

	// ==========================
	// Audit Fields
	// ==========================
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastLogin time.Time `json:"last_login"`
}
