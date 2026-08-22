package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"nigeriaonlinevoting/models"
	"nigeriaonlinevoting/providers"
)

// ============================================================
// Test User Repository
// ============================================================

type testUserRepository struct {
	users map[int]models.User
}

func newTestUserRepository() *testUserRepository {
	return &testUserRepository{
		users: make(map[int]models.User),
	}
}

func (r *testUserRepository) Create(user models.User) error {
	if user.ID <= 0 {
		user.ID = len(r.users) + 1
	}

	r.users[user.ID] = user
	return nil
}

func (r *testUserRepository) GetByEmail(
	email string,
) (*models.User, error) {

	for _, user := range r.users {
		if user.Email == email {
			copy := user
			return &copy, nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *testUserRepository) GetByID(
	id int,
) (*models.User, error) {

	user, ok := r.users[id]

	if !ok {
		return nil, errors.New("user not found")
	}

	copy := user
	return &copy, nil
}

func (r *testUserRepository) Update(
	user models.User,
) error {

	if _, ok := r.users[user.ID]; !ok {
		return errors.New("user not found")
	}

	r.users[user.ID] = user
	return nil
}

func (r *testUserRepository) Delete(id int) error {
	if _, ok := r.users[id]; !ok {
		return errors.New("user not found")
	}

	delete(r.users, id)
	return nil
}

func (r *testUserRepository) GetAll() []models.User {

	users := make([]models.User, 0, len(r.users))

	for _, user := range r.users {
		users = append(users, user)
	}

	return users
}

// ============================================================
// Test Electoral Verification Repository
// ============================================================

type testElectoralVerificationRepository struct {
	records []models.ElectoralVerification
}

func newTestElectoralVerificationRepository() *testElectoralVerificationRepository {
	return &testElectoralVerificationRepository{
		records: []models.ElectoralVerification{},
	}
}

func (r *testElectoralVerificationRepository) Create(
	verification models.ElectoralVerification,
) error {

	verification.ID = len(r.records) + 1

	r.records = append(
		r.records,
		verification,
	)

	return nil
}

func (r *testElectoralVerificationRepository) GetByUserID(
	userID int,
) (*models.ElectoralVerification, error) {

	for _, record := range r.records {

		if record.UserID == userID {
			copy := record
			return &copy, nil
		}
	}

	return nil, errors.New(
		"electoral verification not found",
	)
}

func (r *testElectoralVerificationRepository) GetByID(
	id int,
) (*models.ElectoralVerification, error) {

	for _, record := range r.records {

		if record.ID == id {
			copy := record
			return &copy, nil
		}
	}

	return nil, errors.New(
		"electoral verification not found",
	)
}

func (r *testElectoralVerificationRepository) GetAll() []models.ElectoralVerification {

	return append(
		[]models.ElectoralVerification(nil),
		r.records...,
	)
}

func (r *testElectoralVerificationRepository) Update(
	verification models.ElectoralVerification,
) error {

	for i, record := range r.records {

		if record.ID == verification.ID {

			r.records[i] = verification
			return nil
		}
	}

	return errors.New(
		"electoral verification not found",
	)
}

// ============================================================
// Test Provider
// ============================================================

type testElectoralProvider struct {
	result       providers.VoterVerificationResult
	err          error
	received     providers.VoterVerificationRequest
	contextValue any
}

func (p *testElectoralProvider) VerifyVoter(
	ctx context.Context,
	request providers.VoterVerificationRequest,
) (providers.VoterVerificationResult, error) {

	p.received = request

	p.contextValue = ctx.Value("test-key")

	return p.result, p.err
}

// ============================================================
// Test Helpers
// ============================================================

func createTestUser() models.User {

	return models.User{
		ID:        1,
		FirstName: "Joseph",
		LastName:  "Ochohepo",
		NIN:       "12345678901",
		VIN:       "123456789012345",
		State:     "Benue",
		LGA:       "Agatu",
		DateOfBirth: time.Date(
			1987,
			time.March,
			27,
			0,
			0,
			0,
			0,
			time.UTC,
		),
	}
}

func createTestService(
	user models.User,
	provider providers.ElectoralVerificationProvider,
) (
	*ElectoralVerificationService,
	*testUserRepository,
	*testElectoralVerificationRepository,
) {

	userRepo := newTestUserRepository()

	_ = userRepo.Create(user)

	verificationRepo :=
		newTestElectoralVerificationRepository()

	service :=
		NewElectoralVerificationService(
			verificationRepo,
			userRepo,
			provider,
		)

	return service, userRepo, verificationRepo
}

// ============================================================
// Invalid User ID
// ============================================================

func TestElectoralVerificationInvalidUserID(
	t *testing.T,
) {

	provider := &testElectoralProvider{}

	service, _, _ :=
		createTestService(
			createTestUser(),
			provider,
		)

	_, err := service.VerifyVoter(
		context.Background(),
		0,
	)

	if err == nil {
		t.Fatal(
			"expected invalid user ID error",
		)
	}
}

// ============================================================
// Missing User
// ============================================================

func TestElectoralVerificationUserNotFound(
	t *testing.T,
) {

	provider := &testElectoralProvider{}

	service, _, _ :=
		createTestService(
			createTestUser(),
			provider,
		)

	_, err := service.VerifyVoter(
		context.Background(),
		999,
	)

	if err == nil {
		t.Fatal(
			"expected user not found error",
		)
	}
}

// ============================================================
// Missing First Name
// ============================================================

func TestElectoralVerificationRequiresFirstName(
	t *testing.T,
) {

	user := createTestUser()
	user.FirstName = ""

	provider := &testElectoralProvider{}

	service, _, _ :=
		createTestService(
			user,
			provider,
		)

	_, err := service.VerifyVoter(
		context.Background(),
		user.ID,
	)

	if err == nil {
		t.Fatal(
			"expected first name validation error",
		)
	}
}

// ============================================================
// Missing Last Name
// ============================================================

func TestElectoralVerificationRequiresLastName(
	t *testing.T,
) {

	user := createTestUser()
	user.LastName = ""

	provider := &testElectoralProvider{}

	service, _, _ :=
		createTestService(
			user,
			provider,
		)

	_, err := service.VerifyVoter(
		context.Background(),
		user.ID,
	)

	if err == nil {
		t.Fatal(
			"expected last name validation error",
		)
	}
}

// ============================================================
// Missing VIN
// ============================================================

func TestElectoralVerificationRequiresVIN(
	t *testing.T,
) {

	user := createTestUser()
	user.VIN = ""

	provider := &testElectoralProvider{}

	service, _, _ :=
		createTestService(
			user,
			provider,
		)

	_, err := service.VerifyVoter(
		context.Background(),
		user.ID,
	)

	if err == nil {
		t.Fatal(
			"expected VIN validation error",
		)
	}
}

// ============================================================
// Missing NIN
// ============================================================

func TestElectoralVerificationRequiresNIN(
	t *testing.T,
) {

	user := createTestUser()
	user.NIN = ""

	provider := &testElectoralProvider{}

	service, _, _ :=
		createTestService(
			user,
			provider,
		)

	_, err := service.VerifyVoter(
		context.Background(),
		user.ID,
	)

	if err == nil {
		t.Fatal(
			"expected NIN validation error",
		)
	}
}

// ============================================================
// Successful Verification
// ============================================================

func TestElectoralVerificationSuccess(
	t *testing.T,
) {

	provider := &testElectoralProvider{
		result: providers.VoterVerificationResult{
			Verified:        true,
			RegisteredVoter: true,
			VINMatched:      true,
			IdentityMatched: true,
			ProviderName:    "authorized_test_provider",
			ReferenceID:     "REF-12345",
			Message:         "Voter verified",
		},
	}

	user := createTestUser()

	service, _, repo :=
		createTestService(
			user,
			provider,
		)

	result, err := service.VerifyVoter(
		context.Background(),
		user.ID,
	)

	if err != nil {
		t.Fatalf(
			"expected successful verification, got %v",
			err,
		)
	}

	if result.Status != "verified" {
		t.Fatalf(
			"expected status verified, got %s",
			result.Status,
		)
	}

	if !result.RegisteredVoter {
		t.Fatal(
			"expected RegisteredVoter to be true",
		)
	}

	if !result.VINMatched {
		t.Fatal(
			"expected VINMatched to be true",
		)
	}

	if !result.IdentityMatched {
		t.Fatal(
			"expected IdentityMatched to be true",
		)
	}

	if result.ProviderName != "authorized_test_provider" {
		t.Fatalf(
			"unexpected provider name: %s",
			result.ProviderName,
		)
	}

	if result.ReferenceID != "REF-12345" {
		t.Fatalf(
			"unexpected reference ID: %s",
			result.ReferenceID,
		)
	}

	if len(repo.records) != 1 {
		t.Fatalf(
			"expected one stored record, got %d",
			len(repo.records),
		)
	}
}

// ============================================================
// Provider Does Not Confirm Verification
// ============================================================

func TestElectoralVerificationProviderDoesNotConfirm(
	t *testing.T,
) {

	provider := &testElectoralProvider{
		result: providers.VoterVerificationResult{
			Verified:        false,
			RegisteredVoter: true,
			VINMatched:      true,
			IdentityMatched: true,
			ProviderName:    "authorized_test_provider",
		},
	}

	user := createTestUser()

	service, _, _ :=
		createTestService(
			user,
			provider,
		)

	result, err := service.VerifyVoter(
		context.Background(),
		user.ID,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if result.Status == "verified" {
		t.Fatal(
			"provider did not confirm verification, but result was verified",
		)
	}

	if result.Status != "failed" {
		t.Fatalf(
			"expected failed status, got %s",
			result.Status,
		)
	}
}

// ============================================================
// VIN Mismatch
// ============================================================

func TestElectoralVerificationVINMismatch(
	t *testing.T,
) {

	provider := &testElectoralProvider{
		result: providers.VoterVerificationResult{
			Verified:        false,
			RegisteredVoter: true,
			VINMatched:      false,
			IdentityMatched: true,
			ProviderName:    "authorized_test_provider",
		},
	}

	user := createTestUser()

	service, _, _ :=
		createTestService(
			user,
			provider,
		)

	result, err := service.VerifyVoter(
		context.Background(),
		user.ID,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if result.Status != "failed" {
		t.Fatalf(
			"expected failed status, got %s",
			result.Status,
		)
	}

	if result.VINMatched {
		t.Fatal(
			"expected VIN mismatch",
		)
	}
}

// ============================================================
// Identity Mismatch
// ============================================================

func TestElectoralVerificationIdentityMismatch(
	t *testing.T,
) {

	provider := &testElectoralProvider{
		result: providers.VoterVerificationResult{
			Verified:        false,
			RegisteredVoter: true,
			VINMatched:      true,
			IdentityMatched: false,
			ProviderName:    "authorized_test_provider",
		},
	}

	user := createTestUser()

	service, _, _ :=
		createTestService(
			user,
			provider,
		)

	result, err := service.VerifyVoter(
		context.Background(),
		user.ID,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if result.Status != "failed" {
		t.Fatalf(
			"expected failed status, got %s",
			result.Status,
		)
	}

	if result.IdentityMatched {
		t.Fatal(
			"expected identity mismatch",
		)
	}
}

// ============================================================
// Provider Unavailable
// ============================================================

func TestElectoralVerificationProviderUnavailable(
	t *testing.T,
) {

	provider := &testElectoralProvider{
		result: providers.VoterVerificationResult{
			Verified:        true,
			RegisteredVoter: true,
			VINMatched:      true,
			IdentityMatched: true,
			ProviderName:    "authorized_test_provider",
		},
		err: errors.New(
			"provider unavailable",
		),
	}

	user := createTestUser()

	service, _, repo :=
		createTestService(
			user,
			provider,
		)

	result, err := service.VerifyVoter(
		context.Background(),
		user.ID,
	)

	if err == nil {
		t.Fatal(
			"expected provider error",
		)
	}

	if result.Status != "unavailable" {
		t.Fatalf(
			"expected unavailable status, got %s",
			result.Status,
		)
	}

	// Critical security rule:
	// provider failure must NEVER produce
	// a verified voter.

	if result.RegisteredVoter {
		t.Fatal(
			"provider failure must not produce RegisteredVoter=true",
		)
	}

	if result.VINMatched {
		t.Fatal(
			"provider failure must not produce VINMatched=true",
		)
	}

	if result.IdentityMatched {
		t.Fatal(
			"provider failure must not produce IdentityMatched=true",
		)
	}

	if len(repo.records) != 1 {
		t.Fatal(
			"provider failure should still be recorded for audit purposes",
		)
	}
}

// ============================================================
// Already Verified
// ============================================================

func TestElectoralVerificationAlreadyVerified(
	t *testing.T,
) {

	provider := &testElectoralProvider{}

	user := createTestUser()

	service, _, repo :=
		createTestService(
			user,
			provider,
		)

	_ = repo.Create(
		models.ElectoralVerification{
			UserID: user.ID,
			Status: "verified",
		},
	)

	_, err := service.VerifyVoter(
		context.Background(),
		user.ID,
	)

	if err == nil {
		t.Fatal(
			"expected already verified error",
		)
	}
}

// ============================================================
// Already Pending
// ============================================================

func TestElectoralVerificationAlreadyPending(
	t *testing.T,
) {

	provider := &testElectoralProvider{}

	user := createTestUser()

	service, _, repo :=
		createTestService(
			user,
			provider,
		)

	_ = repo.Create(
		models.ElectoralVerification{
			UserID: user.ID,
			Status: "pending",
		},
	)

	_, err := service.VerifyVoter(
		context.Background(),
		user.ID,
	)

	if err == nil {
		t.Fatal(
			"expected pending verification error",
		)
	}
}

// ============================================================
// Provider Request Data
// ============================================================

func TestElectoralVerificationProviderReceivesCorrectIdentity(
	t *testing.T,
) {

	provider := &testElectoralProvider{
		result: providers.VoterVerificationResult{
			Verified:        false,
			RegisteredVoter: false,
			VINMatched:      false,
			IdentityMatched: false,
			ProviderName:    "authorized_test_provider",
		},
	}

	user := createTestUser()

	service, _, _ :=
		createTestService(
			user,
			provider,
		)

	_, err := service.VerifyVoter(
		context.Background(),
		user.ID,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if provider.received.VIN != user.VIN {
		t.Fatal(
			"provider received incorrect VIN",
		)
	}

	if provider.received.NIN != user.NIN {
		t.Fatal(
			"provider received incorrect NIN",
		)
	}

	if provider.received.FirstName != user.FirstName {
		t.Fatal(
			"provider received incorrect first name",
		)
	}

	if provider.received.LastName != user.LastName {
		t.Fatal(
			"provider received incorrect last name",
		)
	}

	if provider.received.State != user.State {
		t.Fatal(
			"provider received incorrect state",
		)
	}

	if provider.received.LGA != user.LGA {
		t.Fatal(
			"provider received incorrect LGA",
		)
	}

	expectedDOB := "1987-03-27"

	if provider.received.DateOfBirth != expectedDOB {
		t.Fatalf(
			"expected DOB %s, got %s",
			expectedDOB,
			provider.received.DateOfBirth,
		)
	}
}

// ============================================================
// Context Propagation
// ============================================================

func TestElectoralVerificationContextPropagation(
	t *testing.T,
) {

	provider := &testElectoralProvider{
		result: providers.VoterVerificationResult{
			Verified:     false,
			ProviderName: "authorized_test_provider",
		},
	}

	user := createTestUser()

	service, _, _ :=
		createTestService(
			user,
			provider,
		)

	ctx := context.WithValue(
		context.Background(),
		"test-key",
		"test-value",
	)

	_, err := service.VerifyVoter(
		ctx,
		user.ID,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if provider.contextValue != "test-value" {
		t.Fatal(
			"context was not propagated to provider",
		)
	}
}
