package service_test

import (
	"context"
	"testing"

	"github.com/balabanovds/void/internal/domain"
	"github.com/balabanovds/void/internal/models"
	"github.com/balabanovds/void/internal/server/ctxHelper"
	"github.com/balabanovds/void/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestUserService_CreateProfile(t *testing.T) {
	s, u, cl := prepareData(t)
	defer cl()

	newProfile := models.TestNewProfile(t)

	ctx := context.Background()

	// try to create profile not by self and not by admin
	_, err := s.Users().CreateProfile(ctx, newProfile)
	assert.Error(t, err)
	assert.EqualError(t, err, service.ErrNotAllowed.Error())

	// try to create profile by admin - should fail
	ctx = context.WithValue(ctx, ctxHelper.CtxKeyRole, models.Admin)
	_, err = s.Users().CreateProfile(ctx, newProfile)
	assert.Error(t, err)
	assert.EqualError(t, err, service.ErrNotAllowed.Error())

	// create by self - should ok
	ctx = context.WithValue(context.Background(), ctxHelper.CtxKeyEmail, u.Email)
	newProfile.Email = u.Email
	pr, err := s.Users().CreateProfile(ctx, newProfile)
	assert.NoError(t, err)
	assert.NotNil(t, pr.ID)

	// try to create one more
	_, err = s.Users().CreateProfile(ctx, newProfile)
	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrAlreadyExists.Error())

	// create for non existing user
	email := u.Email + "wrong"
	ctx = context.WithValue(context.Background(), ctxHelper.CtxKeyEmail, email)
	newProfile.Email = email
	_, err = s.Users().CreateProfile(ctx, newProfile)
	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrDependencyNotFound.Error())
}

func TestUserService_GetProfile(t *testing.T) {
	s, u, cl := prepareData(t)
	defer cl()

	createdProfile, err := newTestProfile(t, s, u.Email)
	assert.NoError(t, err)

	// try to find non exists profile
	_, err = s.Users().GetProfile(u.Email + "wrong")
	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrNotFound.Error())

	// get exist profile
	gotProfile, err := s.Users().GetProfile(u.Email)
	assert.NoError(t, err)
	assert.Equal(t, createdProfile, gotProfile)
}

func TestUserService_GetAllProfiles(t *testing.T) {
	ts := service.NewTestSuite(t)
	defer ts.Close()

	user1, _ := ts.CreateUser(email, password)
	user2, _ := ts.CreateUser(email+"1", password)

	_, err := newTestProfile(t, ts.Service, user1.Email)
	assert.NoError(t, err)
	_, err = newTestProfile(t, ts.Service, user2.Email)
	assert.NoError(t, err)

	profiles := ts.Service.Users().GetAllProfiles()
	assert.NoError(t, err)
	assert.Len(t, profiles, 2)

}

func TestUserService_UpdateProfile(t *testing.T) {
	s, u, cl := prepareData(t)
	defer cl()

	createdProfile, err := newTestProfile(t, s, u.Email)
	assert.NoError(t, err)

	newFirstName := createdProfile.FirstName + "new"
	newFirstNameRu := createdProfile.Ru.FirstName + "new"

	expectedProfile := createdProfile
	expectedProfile.FirstName = newFirstName
	expectedProfile.Ru.FirstName = newFirstNameRu

	upd := models.UpdateProfile{
		FirstName: newFirstName,
		Ru: models.ProfileRu{
			FirstName: newFirstNameRu,
		},
	}

	cases := []struct {
		name        string
		ctx         context.Context
		expectedErr error
	}{
		{
			name:        "should update by self",
			ctx:         context.WithValue(context.Background(), ctxHelper.CtxKeyEmail, u.Email),
			expectedErr: nil,
		},
		{
			name:        "should fail as non admin and non self",
			ctx:         context.WithValue(context.Background(), ctxHelper.CtxKeyEmail, u.Email+"wrong"),
			expectedErr: service.ErrNotAllowed,
		},
		{
			name:        "should fail while updating by manager",
			ctx:         context.WithValue(context.Background(), ctxHelper.CtxKeyRole, models.Manager),
			expectedErr: service.ErrNotAllowed,
		},
		{
			name:        "should update by admin",
			ctx:         context.WithValue(context.Background(), ctxHelper.CtxKeyRole, models.Admin),
			expectedErr: nil,
		},
	}

	for _, test := range cases {
		testProfile := createdProfile
		err := s.Users().UpdateProfile(test.ctx, &testProfile, upd)
		if test.expectedErr == nil {
			assert.NoError(t, err)
			assert.Equal(t, expectedProfile.FirstName, testProfile.FirstName)
			assert.Equal(t, expectedProfile.LastName, testProfile.LastName)
		} else {
			assert.Error(t, err)
			assert.EqualError(t, err, test.expectedErr.Error())
		}
	}

}

func TestUserService_ChangeRole(t *testing.T) {
	s, u, cl := prepareData(t)
	defer cl()

	createdProfile, err := newTestProfile(t, s, u.Email)
	assert.NoError(t, err)

	testCases := []struct {
		name        string
		ctx         context.Context
		role        models.Role
		expectedErr error
	}{
		{
			name:        "should fail for self",
			ctx:         context.WithValue(context.Background(), ctxHelper.CtxKeyEmail, u.Email),
			role:        models.Manager,
			expectedErr: service.ErrNotAllowed,
		},
		{
			name:        "should fail for manager",
			ctx:         context.WithValue(context.Background(), ctxHelper.CtxKeyRole, models.Manager),
			role:        models.Admin,
			expectedErr: service.ErrNotAllowed,
		},
		{
			name:        "should ok for admin",
			ctx:         context.WithValue(context.Background(), ctxHelper.CtxKeyRole, models.Admin),
			role:        models.Admin,
			expectedErr: nil,
		},
	}

	for _, test := range testCases {
		finalProfile := createdProfile
		err = s.Users().ChangeRole(test.ctx, &finalProfile, test.role)
		if test.expectedErr == nil {
			assert.NoError(t, err)
			assert.Equal(t, test.role, finalProfile.Role)
		} else {
			assert.Error(t, err)
			assert.EqualError(t, err, test.expectedErr.Error())
		}

	}

}

func newTestProfile(t *testing.T, s *service.Service, email string) (models.Profile, error) {
	t.Helper()
	newProfile := models.TestNewProfile(t)
	newProfile.Email = email
	ctx := context.WithValue(context.Background(), ctxHelper.CtxKeyEmail, email)

	return s.Users().CreateProfile(ctx, newProfile)
}
