package providers

import (
	"context"
	"errors"
)

// UnavailableProvider is used when no authorized electoral
// verification service has been configured.
//
// It deliberately does NOT pretend to verify voters.
type UnavailableProvider struct{}

// NewUnavailableProvider creates a provider that refuses
// to perform electoral verification until an authorized
// integration is configured.
func NewUnavailableProvider() *UnavailableProvider {
	return &UnavailableProvider{}
}

// VerifyVoter does not perform a fake verification.
func (p *UnavailableProvider) VerifyVoter(
	ctx context.Context,
	request VoterVerificationRequest,
) (VoterVerificationResult, error) {

	return VoterVerificationResult{
			Verified:        false,
			RegisteredVoter: false,
			VINMatched:      false,
			IdentityMatched: false,
			ProviderName:    "unavailable",
			Message:         "No authorized electoral verification provider is configured.",
		},
		errors.New(
			"electoral verification service is unavailable",
		)
}