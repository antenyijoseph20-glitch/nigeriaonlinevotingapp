package services

import (
	"errors"
	"testing"

	"nigeriaonlinevoting/models"
)

// =====================================
// Fake User Repository
// =====================================

type fakeUserRepository struct {
	users []models.User
	err   error
}

func (f *fakeUserRepository) Create(
	user models.User,
) error {
	f.users = append(f.users, user)
	return nil
}

func (f *fakeUserRepository) Update(
	user models.User,
) error {
	for i := range f.users {
		if f.users[i].ID == user.ID {
			f.users[i] = user
			return nil
		}
	}

	return errors.New("user not found")
}

func (f *fakeUserRepository) Delete(
	id int,
) error {
	for i := range f.users {
		if f.users[i].ID == id {
			f.users = append(
				f.users[:i],
				f.users[i+1:]...,
			)
			return nil
		}
	}

	return errors.New("user not found")
}

func (f *fakeUserRepository) GetByID(
	id int,
) (*models.User, error) {

	if f.err != nil {
		return nil, f.err
	}

	for i := range f.users {
		if f.users[i].ID == id {
			return &f.users[i], nil
		}
	}

	return nil, errors.New("user not found")
}

func (f *fakeUserRepository) GetByEmail(
	email string,
) (*models.User, error) {

	for i := range f.users {
		if f.users[i].Email == email {
			return &f.users[i], nil
		}
	}

	return nil, errors.New("user not found")
}

func (f *fakeUserRepository) GetAll() []models.User {
	return f.users
}

// =====================================
// Fake Voter Repository
// =====================================

type fakeVoterRepository struct {
	voters []models.Voter
	err    error
}

func (f *fakeVoterRepository) Create(
	voter models.Voter,
) error {

	if f.err != nil {
		return f.err
	}

	voter.ID = len(f.voters) + 1

	f.voters = append(
		f.voters,
		voter,
	)

	return nil
}

func (f *fakeVoterRepository) Update(
	voter models.Voter,
) error {

	if f.err != nil {
		return f.err
	}

	for i := range f.voters {

		if f.voters[i].ID == voter.ID {

			f.voters[i] = voter

			return nil
		}
	}

	return errors.New("voter not found")
}

func (f *fakeVoterRepository) Delete(
	id int,
) error {

	if f.err != nil {
		return f.err
	}

	for i := range f.voters {

		if f.voters[i].ID == id {

			f.voters = append(
				f.voters[:i],
				f.voters[i+1:]...,
			)

			return nil
		}
	}

	return errors.New("voter not found")
}

func (f *fakeVoterRepository) GetByID(
	id int,
) (*models.Voter, error) {

	if f.err != nil {
		return nil, f.err
	}

	for i := range f.voters {

		if f.voters[i].ID == id {
			return &f.voters[i], nil
		}
	}

	return nil, errors.New("voter not found")
}

func (f *fakeVoterRepository) GetByUserID(
	userID int,
) (*models.Voter, error) {

	if f.err != nil {
		return nil, f.err
	}

	for i := range f.voters {

		if f.voters[i].UserID == userID {
			return &f.voters[i], nil
		}
	}

	return nil, errors.New("voter not found")
}

func (f *fakeVoterRepository) GetByNIN(
	nin string,
) (*models.Voter, error) {

	if f.err != nil {
		return nil, f.err
	}

	for i := range f.voters {

		if f.voters[i].NIN == nin {
			return &f.voters[i], nil
		}
	}

	return nil, errors.New("voter not found")
}

func (f *fakeVoterRepository) GetByVIN(
	vin string,
) (*models.Voter, error) {

	if f.err != nil {
		return nil, f.err
	}

	for i := range f.voters {

		if f.voters[i].VIN == vin {
			return &f.voters[i], nil
		}
	}

	return nil, errors.New("voter not found")
}

func (f *fakeVoterRepository) GetByPVCNumber(
	pvcNumber string,
) (*models.Voter, error) {

	if f.err != nil {
		return nil, f.err
	}

	for i := range f.voters {

		if f.voters[i].PVCNumber == pvcNumber {
			return &f.voters[i], nil
		}
	}

	return nil, errors.New("voter not found")
}

func (f *fakeVoterRepository) GetAll() []models.Voter {
	return f.voters
}

// =====================================
// Test Data Helpers
// =====================================

func testUser() models.User {

	return models.User{
		ID:            1,
		FirstName:     "Joseph",
		LastName:      "Ochohepo",
		Email:         "joseph@example.com",
		AccountActive: true,
	}
}

func testVoter() models.Voter {

	return models.Voter{
		UserID:      1,
		PVCNumber:   "PVC001",
		NIN:         "NIN001",
		VIN:         "VIN001",
		State:       "Benue",
		LGA:         "Otukpo",
		Ward:        "Ward 1",
		PollingUnit: "PU001",
	}
}

func newTestVoterService() (
	*VoterService,
	*fakeVoterRepository,
	*fakeUserRepository,
) {

	voterRepo := &fakeVoterRepository{}

	userRepo := &fakeUserRepository{
		users: []models.User{
			testUser(),
		},
	}

	service := NewVoterService(
		voterRepo,
		userRepo,
	)

	return service, voterRepo, userRepo
}

// =====================================
// Register Voter Tests
// =====================================

func TestRegisterVoterSuccess(t *testing.T) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()

	err := service.RegisterVoter(voter)

	if err != nil {
		t.Fatalf(
			"expected registration to succeed, got %v",
			err,
		)
	}

	if len(voterRepo.voters) != 1 {
		t.Fatalf(
			"expected 1 voter, got %d",
			len(voterRepo.voters),
		)
	}

	saved := voterRepo.voters[0]

	if !saved.IsRegistered {
		t.Error("expected voter to be registered")
	}

	if saved.IsVerified {
		t.Error("new voter should not be verified")
	}

	if saved.IsEligible {
		t.Error("new voter should not be eligible")
	}

	if saved.HasVoted {
		t.Error("new voter should not have voted")
	}

	if saved.ID <= 0 {
		t.Error("expected voter ID to be generated")
	}
}

func TestRegisterVoterRejectsInvalidUserID(t *testing.T) {

	service, _, _ := newTestVoterService()

	voter := testVoter()
	voter.UserID = 0

	err := service.RegisterVoter(voter)

	if err == nil {
		t.Fatal("expected invalid user ID to be rejected")
	}
}

func TestRegisterVoterRejectsEmptyNIN(t *testing.T) {

	service, _, _ := newTestVoterService()

	voter := testVoter()
	voter.NIN = ""

	err := service.RegisterVoter(voter)

	if err == nil {
		t.Fatal("expected empty NIN to be rejected")
	}
}

func TestRegisterVoterRejectsEmptyVIN(t *testing.T) {

	service, _, _ := newTestVoterService()

	voter := testVoter()
	voter.VIN = ""

	err := service.RegisterVoter(voter)

	if err == nil {
		t.Fatal("expected empty VIN to be rejected")
	}
}

func TestRegisterVoterRejectsEmptyPVC(t *testing.T) {

	service, _, _ := newTestVoterService()

	voter := testVoter()
	voter.PVCNumber = ""

	err := service.RegisterVoter(voter)

	if err == nil {
		t.Fatal("expected empty PVC number to be rejected")
	}
}

func TestRegisterVoterRejectsMissingUser(t *testing.T) {

	service, _, userRepo := newTestVoterService()

	userRepo.users = nil

	err := service.RegisterVoter(testVoter())

	if err == nil {
		t.Fatal("expected missing user to be rejected")
	}
}

func TestRegisterVoterRejectsDuplicateUser(t *testing.T) {

	service, _, _ := newTestVoterService()

	if err := service.RegisterVoter(testVoter()); err != nil {
		t.Fatalf(
			"first registration failed: %v",
			err,
		)
	}

	err := service.RegisterVoter(testVoter())

	if err == nil {
		t.Fatal(
			"expected duplicate user registration to be rejected",
		)
	}
}

func TestRegisterVoterRejectsDuplicateNIN(t *testing.T) {

	service, _, _ := newTestVoterService()

	first := testVoter()

	if err := service.RegisterVoter(first); err != nil {
		t.Fatalf(
			"first registration failed: %v",
			err,
		)
	}

	second := testVoter()
	second.UserID = 2
	second.VIN = "VIN002"
	second.PVCNumber = "PVC002"

	err := service.RegisterVoter(second)

	if err == nil {
		t.Fatal(
			"expected duplicate NIN to be rejected",
		)
	}
}

func TestRegisterVoterRejectsDuplicateVIN(t *testing.T) {

	service, _, userRepo := newTestVoterService()

	userRepo.users = append(
		userRepo.users,
		models.User{
			ID:            2,
			FirstName:     "Test",
			LastName:      "User",
			Email:         "test@example.com",
			AccountActive: true,
		},
	)

	if err := service.RegisterVoter(testVoter()); err != nil {
		t.Fatalf(
			"first registration failed: %v",
			err,
		)
	}

	second := testVoter()
	second.UserID = 2
	second.NIN = "NIN002"
	second.PVCNumber = "PVC002"

	err := service.RegisterVoter(second)

	if err == nil {
		t.Fatal(
			"expected duplicate VIN to be rejected",
		)
	}
}

func TestRegisterVoterRejectsDuplicatePVC(t *testing.T) {

	service, _, userRepo := newTestVoterService()

	userRepo.users = append(
		userRepo.users,
		models.User{
			ID:            2,
			FirstName:     "Test",
			LastName:      "User",
			Email:         "test@example.com",
			AccountActive: true,
		},
	)

	if err := service.RegisterVoter(testVoter()); err != nil {
		t.Fatalf(
			"first registration failed: %v",
			err,
		)
	}

	second := testVoter()
	second.UserID = 2
	second.NIN = "NIN002"
	second.VIN = "VIN002"

	err := service.RegisterVoter(second)

	if err == nil {
		t.Fatal(
			"expected duplicate PVC number to be rejected",
		)
	}
}

// =====================================
// Get Voter Tests
// =====================================

func TestGetVoterRejectsInvalidID(t *testing.T) {

	service, _, _ := newTestVoterService()

	_, err := service.GetVoter(0)

	if err == nil {
		t.Fatal("expected invalid voter ID to be rejected")
	}
}

func TestGetVoterSuccess(t *testing.T) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	result, err := service.GetVoter(1)

	if err != nil {
		t.Fatalf(
			"expected voter lookup to succeed: %v",
			err,
		)
	}

	if result.ID != 1 {
		t.Fatalf(
			"expected voter ID 1, got %d",
			result.ID,
		)
	}
}

// =====================================
// Verify Voter Tests
// =====================================

func TestVerifyVoterSuccess(t *testing.T) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1
	voter.IsRegistered = true

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	err := service.VerifyVoter(1)

	if err != nil {
		t.Fatalf(
			"expected verification to succeed: %v",
			err,
		)
	}

	verified := voterRepo.voters[0]

	if !verified.IsVerified {
		t.Error("expected voter to be verified")
	}

	if !verified.IsEligible {
		t.Error("expected verified voter to become eligible")
	}
}

func TestVerifyVoterRejectsSuspendedVoter(t *testing.T) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1
	voter.IsRegistered = true
	voter.IsSuspended = true

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	err := service.VerifyVoter(1)

	if err == nil {
		t.Fatal(
			"expected suspended voter verification to fail",
		)
	}
}

func TestVerifyVoterRejectsLockedVoter(t *testing.T) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1
	voter.IsRegistered = true
	voter.AccountLocked = true

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	err := service.VerifyVoter(1)

	if err == nil {
		t.Fatal(
			"expected locked voter verification to fail",
		)
	}
}

// =====================================
// Suspend / Reinstate Tests
// =====================================

func TestSuspendVoterSuccess(t *testing.T) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1
	voter.IsRegistered = true
	voter.IsVerified = true
	voter.IsEligible = true

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	err := service.SuspendVoter(1)

	if err != nil {
		t.Fatalf(
			"expected suspension to succeed: %v",
			err,
		)
	}

	if !voterRepo.voters[0].IsSuspended {
		t.Error("expected voter to be suspended")
	}

	if voterRepo.voters[0].IsEligible {
		t.Error(
			"suspended voter must not remain eligible",
		)
	}
}

func TestSuspendVoterRejectsAlreadySuspended(
	t *testing.T,
) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1
	voter.IsSuspended = true

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	err := service.SuspendVoter(1)

	if err == nil {
		t.Fatal(
			"expected already suspended voter to be rejected",
		)
	}
}

func TestReinstateVoterSuccess(t *testing.T) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1
	voter.IsRegistered = true
	voter.IsVerified = true
	voter.IsSuspended = true
	voter.IsEligible = false

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	err := service.ReinstateVoter(1)

	if err != nil {
		t.Fatalf(
			"expected reinstatement to succeed: %v",
			err,
		)
	}

	if voterRepo.voters[0].IsSuspended {
		t.Error(
			"expected voter to no longer be suspended",
		)
	}

	if !voterRepo.voters[0].IsEligible {
		t.Error(
			"expected verified registered voter to become eligible",
		)
	}
}

func TestReinstateVoterRejectsActiveVoter(t *testing.T) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1
	voter.IsRegistered = true
	voter.IsVerified = true
	voter.IsSuspended = false

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	err := service.ReinstateVoter(1)

	if err == nil {
		t.Fatal(
			"expected active voter reinstatement to fail",
		)
	}
}

// =====================================
// Biometric Tests
// =====================================

func TestEnrollFaceSuccess(t *testing.T) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	err := service.EnrollFace(1)

	if err != nil {
		t.Fatalf(
			"expected face enrollment to succeed: %v",
			err,
		)
	}

	if !voterRepo.voters[0].FaceEnrolled {
		t.Error("expected face to be enrolled")
	}
}

func TestVerifyFaceRequiresEnrollment(t *testing.T) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	err := service.VerifyFace(1)

	if err == nil {
		t.Fatal(
			"expected face verification without enrollment to fail",
		)
	}
}

func TestVerifyFaceSuccess(t *testing.T) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1
	voter.FaceEnrolled = true

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	err := service.VerifyFace(1)

	if err != nil {
		t.Fatalf(
			"expected face verification to succeed: %v",
			err,
		)
	}

	if !voterRepo.voters[0].FaceVerified {
		t.Error("expected face to be verified")
	}
}

func TestEnrollFingerprintSuccess(t *testing.T) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	err := service.EnrollFingerprint(1)

	if err != nil {
		t.Fatalf(
			"expected fingerprint enrollment to succeed: %v",
			err,
		)
	}

	if !voterRepo.voters[0].FingerprintEnrolled {
		t.Error(
			"expected fingerprint to be enrolled",
		)
	}
}

func TestVerifyFingerprintRequiresEnrollment(
	t *testing.T,
) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	err := service.VerifyFingerprint(1)

	if err == nil {
		t.Fatal(
			"expected fingerprint verification without enrollment to fail",
		)
	}
}

func TestVerifyFingerprintSuccess(t *testing.T) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1
	voter.FingerprintEnrolled = true

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	err := service.VerifyFingerprint(1)

	if err != nil {
		t.Fatalf(
			"expected fingerprint verification to succeed: %v",
			err,
		)
	}

	if !voterRepo.voters[0].FingerprintVerified {
		t.Error(
			"expected fingerprint to be verified",
		)
	}
}

// =====================================
// Voting Status Tests
// =====================================

func TestMarkAsVotedSuccess(t *testing.T) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1
	voter.IsRegistered = true
	voter.IsVerified = true
	voter.IsEligible = true

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	err := service.MarkAsVoted(1)

	if err != nil {
		t.Fatalf(
			"expected voter to be marked as voted: %v",
			err,
		)
	}

	if !voterRepo.voters[0].HasVoted {
		t.Error(
			"expected HasVoted to be true",
		)
	}
}

func TestMarkAsVotedRejectsSuspendedVoter(t *testing.T) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1
	voter.IsEligible = true
	voter.IsSuspended = true

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	err := service.MarkAsVoted(1)

	if err == nil {
		t.Fatal(
			"expected suspended voter to be rejected",
		)
	}
}

func TestMarkAsVotedRejectsIneligibleVoter(t *testing.T) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1
	voter.IsEligible = false

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	err := service.MarkAsVoted(1)

	if err == nil {
		t.Fatal(
			"expected ineligible voter to be rejected",
		)
	}
}

func TestMarkAsVotedRejectsSecondVote(t *testing.T) {

	service, voterRepo, _ := newTestVoterService()

	voter := testVoter()
	voter.ID = 1
	voter.IsEligible = true
	voter.HasVoted = true

	voterRepo.voters = append(
		voterRepo.voters,
		voter,
	)

	err := service.MarkAsVoted(1)

	if err == nil {
		t.Fatal(
			"expected second vote to be rejected",
		)
	}
}

// =====================================
// Repository Error Tests
// =====================================

func TestRegisterVoterPropagatesUserRepositoryError(
	t *testing.T,
) {

	service, _, userRepo := newTestVoterService()

	userRepo.err = errors.New(
		"user repository failure",
	)

	err := service.RegisterVoter(testVoter())

	if err == nil {
		t.Fatal(
			"expected repository error to be returned",
		)
	}
}

func TestGetVoterPropagatesRepositoryError(
	t *testing.T,
) {

	service, _, voterRepo := newTestVoterService()

	voterRepo.err = errors.New(
		"voter repository failure",
	)

	_, err := service.GetVoter(1)

	if err == nil {
		t.Fatal(
			"expected repository error to be returned",
		)
	}
}
