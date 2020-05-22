package service_test

import (
	"context"
	"github.com/balabanovds/void/internal/domain"
	"github.com/balabanovds/void/internal/models"
	"github.com/balabanovds/void/internal/service"
	"testing"

	"github.com/balabanovds/void/internal/server/ctxHelper"

	"github.com/stretchr/testify/assert"
)

const (
	email           string = "test"
	password        string = "pass"
	passwordConfirm string = "pass"
)

func prepareData(t *testing.T) (*service.Service, models.User, func()) {
	t.Helper()

	ts := service.NewTestSuite(t)
	user, err := ts.CreateUser(email, password)
	if err != nil {
		t.Fatal(err)
	}

	return ts.Service, user, ts.Close
}

func TestUserService_Create(t *testing.T) {
	ts := service.NewTestSuite(t)
	defer ts.Close()

	_, err := ts.Service.Users().Create(email, password, passwordConfirm+"somestuff")
	assert.Error(t, err)
	assert.EqualError(t, err, service.ErrPasswdNotMatch.Error())

	u, err := ts.Service.Users().Create(email, password, passwordConfirm)

	assert.NoError(t, err)
	assert.Equal(t, email, u.Email)
}

func TestUserService_Authenticate(t *testing.T) {
	s, user, cl := prepareData(t)
	defer cl()

	authUser, err := s.Users().Authenticate(email, password)

	assert.NoError(t, err)
	assert.Equal(t, user, authUser)

	_, err = s.Users().Authenticate(email, password+"not_match")

	assert.Error(t, err)
	assert.EqualError(t, err, service.ErrFailedAuthenticate.Error())
}

func TestUserService_GetByEmail(t *testing.T) {
	s, user, cl := prepareData(t)
	defer cl()

	_, err := s.Users().GetByEmail(email + "somestuff")
	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrNotFound.Error())

	gotUser, err := s.Users().GetByEmail(email)

	assert.NoError(t, err)
	assert.Equal(t, user, gotUser)
}

func TestUserService_UpdatePassword(t *testing.T) {
	s, user, cl := prepareData(t)
	defer cl()

	ctx := context.Background()
	newPass := "new__password"

	// try with empty context - should fail
	err := s.Users().UpdatePassword(ctx, user.Email, newPass)
	assert.Error(t, err)
	assert.EqualError(t, err, service.ErrNotAllowed.Error())

	// try as admin but not self - should fail
	ctx = context.WithValue(ctx, ctxHelper.CtxKeyIsAdmin, true)
	err = s.Users().UpdatePassword(ctx, user.Email, newPass)
	assert.Error(t, err)
	assert.EqualError(t, err, service.ErrNotAllowed.Error())

	// try as self - should be ok
	ctx = context.WithValue(context.Background(), ctxHelper.CtxKeyEmail, user.Email)
	err = s.Users().UpdatePassword(ctx, user.Email, newPass)
	assert.NoError(t, err)

	// check if password has been changed
	authUser, err := s.Users().Authenticate(user.Email, newPass)
	assert.NoError(t, err)
	assert.NotEqual(t, user.HashedPassword, authUser.HashedPassword)
}

func TestUserService_ToggleActive(t *testing.T) {
	s, user, cl := prepareData(t)
	defer cl()

	// try to toggle active by self - should fail
	ctx := context.WithValue(context.Background(), ctxHelper.CtxKeyEmail, user.Email)
	err := s.Users().ToggleActive(ctx, user.Email)
	assert.Error(t, err)
	assert.EqualError(t, err, service.ErrNotAllowed.Error())

	// try to toggle active by admin - should be ok
	ctx = context.WithValue(ctx, ctxHelper.CtxKeyIsAdmin, true)
	err = s.Users().ToggleActive(ctx, user.Email)
	assert.NoError(t, err)

	updatedUser, err := s.Users().GetByEmail(user.Email)
	assert.NoError(t, err)
	assert.Equal(t, !user.Active, updatedUser.Active)
}

func TestUserService_Delete(t *testing.T) {
	s, user, cl := prepareData(t)
	defer cl()

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxHelper.CtxKeyEmail, user.Email+"wrong")

	// try to delete not self and not admin - should fail
	err := s.Users().Delete(ctx, user.Email)
	assert.Error(t, err)
	assert.EqualError(t, err, service.ErrNotAllowed.Error())

	// try to delete by admin - should be ok
	ctx = context.WithValue(context.Background(), ctxHelper.CtxKeyIsAdmin, true)
	err = s.Users().Delete(ctx, user.Email)
	assert.NoError(t, err)

	_, err = s.Users().GetByEmail(user.Email)
	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrNotFound.Error())

	// create again and try to delete by self
	_, err = s.Users().Create(email, password, passwordConfirm)
	assert.NoError(t, err)
	ctx = context.WithValue(context.Background(), ctxHelper.CtxKeyEmail, user.Email)
	err = s.Users().Delete(ctx, user.Email)
	assert.NoError(t, err)

	_, err = s.Users().GetByEmail(user.Email)
	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrNotFound.Error())

}
