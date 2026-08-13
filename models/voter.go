package models

import "time"

// Voter represents a person who is registered
// and eligible to participate in an election.
//
// User represents the application account.
// Voter represents the electoral record.
type Voter struct {

	// =====================================
	// Primary Identifier
	// =====================================

	ID int `json:"id"`

	// =====================================
	// Linked User
	// =====================================

	// User account associated with this voter.
	UserID int `json:"user_id"`

	// =====================================
	// Electoral Identification
	// =====================================

	// Permanent Voter Card number.
	PVCNumber string `json:"pvc_number"`

	// National Identification Number.
	NIN string `json:"nin"`

	// Voter Identification Number.
	VIN string `json:"vin"`

	// =====================================
	// Electoral Location
	// =====================================

	State string `json:"state"`

	LGA string `json:"lga"`

	Ward string `json:"ward"`

	PollingUnit string `json:"polling_unit"`

	PollingCode string `json:"polling_code"`

	RegistrationArea string `json:"registration_area"`

	// =====================================
	// Registration Status
	// =====================================

	// Indicates that the voter has been
	// registered in the electoral system.
	IsRegistered bool `json:"is_registered"`

	// Indicates that the voter's identity
	// has been verified.
	IsVerified bool `json:"is_verified"`

	// Indicates that the voter is currently
	// eligible to vote.
	IsEligible bool `json:"is_eligible"`

	// Indicates that the voter has been
	// suspended from voting.
	IsSuspended bool `json:"is_suspended"`

	// =====================================
	// Biometric Verification
	// =====================================

	FaceEnrolled bool `json:"face_enrolled"`

	FingerprintEnrolled bool `json:"fingerprint_enrolled"`

	FaceVerified bool `json:"face_verified"`

	FingerprintVerified bool `json:"fingerprint_verified"`

	LastVerification time.Time `json:"last_verification"`

	// =====================================
	// Voting Status
	// =====================================

	// This is the overall account-level
	// voting status.
	//
	// Individual election/position duplicate
	// voting will be handled by Vote records.
	HasVoted bool `json:"has_voted"`

	LastVoteTime time.Time `json:"last_vote_time"`

	// =====================================
	// Security
	// =====================================

	FailedLoginAttempts int `json:"failed_login_attempts"`

	AccountLocked bool `json:"account_locked"`

	// =====================================
	// Audit Fields
	// =====================================

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}
