package providers

import "context"

// VoterVerificationRequest contains the minimum information
// required to ask an authorized electoral verification provider
// to verify a voter.
//
// This structure intentionally does not assume that a particular
// electoral authority exposes a public API.
type VoterVerificationRequest struct {
	VIN         string
	NIN         string
	FirstName   string
	LastName    string
	DateOfBirth string
	State       string
	LGA         string
}

// VoterVerificationResult represents the result returned by
// an authorized electoral verification provider.
type VoterVerificationResult struct {
	Verified       bool
	RegisteredVoter bool
	VINMatched     bool
	IdentityMatched bool

	ProviderName string
	ReferenceID  string
	Message      string
}

// ElectoralVerificationProvider defines the boundary between
// our application and an external electoral verification system.
//
// A real implementation must only be connected to an authorized
// and legally permitted data source.
type ElectoralVerificationProvider interface {
	VerifyVoter(
		ctx context.Context,
		request VoterVerificationRequest,
	) (VoterVerificationResult, error)
}