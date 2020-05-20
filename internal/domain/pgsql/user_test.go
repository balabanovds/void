package pgsql_test

import (
	"testing"

	"github.com/balabanovds/void/internal/domain"
	"github.com/balabanovds/void/internal/domain/pgsql"
	"github.com/balabanovds/void/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo_Create(t *testing.T) {
	u := models.TestUser(t)
	s, cleanup := pgsql.TestDB(t, databaseUrl)
	defer cleanup("users")

	createdUser, err := s.Users().Create(u.Email, u.HashedPassword)

	assert.NoError(t, err)
	assert.NotZero(t, createdUser.ID)
	assert.True(t, createdUser.Active)
}

func TestUserRepo_CreateDuplicatedEmail(t *testing.T) {
	u := models.TestUser(t)
	s, cleanup := pgsql.TestDB(t, databaseUrl)
	defer cleanup("users")

	createdUser, err := s.Users().Create(u.Email, u.HashedPassword)

	assert.NoError(t, err)
	assert.NotZero(t, createdUser.ID)

	_, err = s.Users().Create(u.Email, "p")

	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrDuplicateEmail.Error())

}

func TestUserRepo_Get(t *testing.T) {
	s, u, cl := getTestData(t)
	defer cl()

	gotUser, err := s.Users().Get(u.ID)

	assert.NoError(t, err)
	assert.Equal(t, u.ID, gotUser.ID)
	assert.Equal(t, u.Active, gotUser.Active)
}

func TestUserRepo_GetNotFound(t *testing.T) {
	s, cleanup := pgsql.TestDB(t, databaseUrl)
	defer cleanup("users")

	_, err := s.Users().Get(1)

	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrNotFound.Error())
}

func TestUserRepo_GetByEmail(t *testing.T) {
	s, u, cl := getTestData(t)
	defer cl()

	gotUser, err := s.Users().GetByEmail(u.Email)

	assert.NoError(t, err)
	assert.Equal(t, u.ID, gotUser.ID)
	assert.Equal(t, u.Active, gotUser.Active)
}

func TestUserRepo_GetByEmailNotFound(t *testing.T) {
	s, cleanup := pgsql.TestDB(t, databaseUrl)
	defer cleanup("users")

	_, err := s.Users().GetByEmail("1")

	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrNotFound.Error())
}

func TestUserRepo_Update(t *testing.T) {
	s, user, cl := getTestData(t)
	defer cl()

	newActive := !user.Active
	newPass := "newPass"

	err := s.Users().Update(&user, newPass, newActive)

	assert.NoError(t, err)
	assert.Equal(t, newPass, user.HashedPassword)
	assert.Equal(t, newActive, user.Active)

	gotUser, err := s.Users().GetByEmail(user.Email)

	assert.NoError(t, err)
	assert.Equal(t, newActive, gotUser.Active)
	assert.Equal(t, newPass, gotUser.HashedPassword)
}

func TestUserRepo_Delete(t *testing.T) {
	s, u, cl := getTestData(t)
	defer cl()

	s.Users().Delete(u.ID)

	_, err := s.Users().Get(u.ID)

	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrNotFound.Error())
}

func getTestData(t *testing.T) (domain.Storage, models.User, func()) {
	u := models.TestUser(t)
	s, cleanup := pgsql.TestDB(t, databaseUrl)

	user, err := s.Users().Create(u.Email, u.HashedPassword)
	if err != nil {
		t.Fatal(err)
	}

	return s, user, func() {
		cleanup("users")
	}
}
