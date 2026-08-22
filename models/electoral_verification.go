package models

import "time"

// ElectoralVerification records the result of checking a voter
// against an authorized electoral verification provider.
//
// This record is deliberately separate from the application's
// administrative verification process.
//
// An administrator approving a submitted form does NOT by itself
// mean that the voter has been confirmed against an official
// electoral register.
type ElectoralVerification struct {
	// Primary key.
	ID int `json:"id"`

	// Application user associated with this verification.
	UserID int `json:"user_id"`

	// Verification status.
	//
	// Possible values:
	// pending
	// verified
	// failed
	// unavailable
	Status string `json:"status"`

	// Provider used for the verification.
	//
	// Examples:
	// authorized_provider
	// unavailable
	ProviderName string `json:"provider_name"`

	// Reference returned by the authorized provider.
	//
	// Do not place NIN, VIN, passwords, or other sensitive
	// personal information in this field.
	ReferenceID string `json:"reference_id"`

	// Whether the provider confirmed that the voter
	// is registered.
	RegisteredVoter bool `json:"registered_voter"`

	// Whether the supplied VIN matched the provider's
	// electoral record.
	VINMatched bool `json:"vin_matched"`

	// Whether the supplied identity information matched
	// the provider's record.
	IdentityMatched bool `json:"identity_matched"`

	// Human-readable provider response.
	Message string `json:"message"`

	// When verification was requested.
	RequestedAt time.Time `json:"requested_at"`

	// When the provider returned a result.
	VerifiedAt time.Time `json:"verified_at"`
}