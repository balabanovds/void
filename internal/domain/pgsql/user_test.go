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

	user, err := ts.Storage.Users().Create(ts.User.Email, ts.User.HashedPassword)
	if err != nil {
		t.Fatal(err)
	}

	return ts.Storage, user, ts.Close
}

func TestUserRepo_Create(t *testing.T) {
	ts := NewTestSuite(t)
	defer ts.Close()

	createdUser, err := ts.Storage.Users().Create(ts.User.Email, ts.User.HashedPassword)

	assert.NoError(t, err)
	assert.NotZero(t, createdUser.ID)
	assert.True(t, createdUser.Active)
}

func TestUserRepo_CreateDuplicatedEmail(t *testing.T) {
	ts := NewTestSuite(t)
	defer ts.Close()

	_, _ = ts.Storage.Users().Create(ts.User.Email, ts.User.HashedPassword)
	_, err := ts.Storage.Users().Create(ts.User.Email, []byte("p"))

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
	s, cleanup := TestDB(t)
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
