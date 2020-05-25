package pgsql

import (
	"testing"

	"github.com/balabanovds/void/internal/domain"
	"github.com/balabanovds/void/internal/models"
	"github.com/stretchr/testify/assert"
)

var (
	ts TestSuite
)

func prepareData(t *testing.T) (domain.Storage, models.User, func()) {
	ts = NewTestSuite(t)

	user := ts.CreateDefaultUser(t)

	return ts.Storage, user, ts.Close
}

func TestUserRepo_Create(t *testing.T) {
	ts := NewTestSuite(t)
	defer ts.Close()

	newUser := models.NewUser{
		Email:          ts.User.Email,
		HashedPassword: ts.User.HashedPassword,
	}

	createdUser, err := ts.Storage.Users().Create(newUser)

	assert.NoError(t, err)
	assert.NotZero(t, createdUser.ID)
	assert.True(t, createdUser.Active)
}

func TestUserRepo_CreateDuplicatedEmail(t *testing.T) {
	ts := NewTestSuite(t)
	defer ts.Close()

	newUser := models.NewUser{
		Email:          ts.User.Email,
		HashedPassword: ts.User.HashedPassword,
	}

	_, _ = ts.Storage.Users().Create(newUser)
	newUser.HashedPassword = []byte("p")
	_, err := ts.Storage.Users().Create(newUser)

	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrDuplicateEmail.Error())
}

func TestUserRepo_GetByEmail(t *testing.T) {
	s, u, cl := prepareData(t)
	defer cl()

	gotUser, err := s.Users().Get(u.Email)

	assert.NoError(t, err)
	assert.Equal(t, u.ID, gotUser.ID)
	assert.Equal(t, u.Active, gotUser.Active)
}

func TestUserRepo_GetByEmailNotFound(t *testing.T) {
	s, cleanup := testDB(t)
	defer cleanup("users")

	_, err := s.Users().Get("1")

	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrNotFound.Error())
}

func TestUserRepo_Update(t *testing.T) {
	s, user, cl := prepareData(t)
	defer cl()

	newActive := !user.Active
	newPass := []byte("newPass")

	err := s.Users().Update(&user, newPass, newActive)

	assert.NoError(t, err)
	assert.Equal(t, newPass, user.HashedPassword)
	assert.Equal(t, newActive, user.Active)

	gotUser, err := s.Users().Get(user.Email)

	assert.NoError(t, err)
	assert.Equal(t, newActive, gotUser.Active)
	assert.Equal(t, newPass, gotUser.HashedPassword)
}

func TestUserRepo_Delete(t *testing.T) {
	s, u, cl := prepareData(t)
	defer cl()

	s.Users().Delete(u.Email)

	_, err := s.Users().Get(u.Email)

	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrNotFound.Error())
}
