package pgsql

import (
	"github.com/balabanovds/void/internal/domain"
	"github.com/balabanovds/void/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProfileRepo_Create(t *testing.T) {
	ts := NewTestSuite(t)
	defer ts.Close()

	newUser := models.NewUser{
		Email: ts.User.Email,
		HashedPassword: ts.User.HashedPassword,
	}

	u, err := ts.Storage.Users().Create(newUser)
	assert.NoError(t, err)

	newProfile := models.TestNewProfile(t)

	newProfile.Email = u.Email+"wrong"
	_, err = ts.Storage.Profiles().Create(newProfile)
	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrDependencyNotFound.Error())

	newProfile.Email = u.Email

	createdProfile, err := ts.Storage.Profiles().Create(newProfile)
	assert.NoError(t, err)
	assert.NotNil(t, createdProfile.ID)

	_, err = ts.Storage.Profiles().Create(newProfile)
	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrAlreadyExists.Error())

}

func TestProfileRepo_Get(t *testing.T) {
	ts := NewTestSuite(t)
	defer ts.Close()

	newUser := models.NewUser{
		Email: ts.User.Email,
		HashedPassword: ts.User.HashedPassword,
	}

	u, err := ts.Storage.Users().Create(newUser)
	assert.NoError(t, err)

	newProfile := models.TestNewProfile(t)
	newProfile.Email = u.Email
	createdProfile, err := ts.Storage.Profiles().Create(newProfile)
	assert.NoError(t, err)

	_, err = ts.Storage.Profiles().Get(u.Email+"wrong")
	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrNotFound.Error())

	foundProfile, err := ts.Storage.Profiles().Get(u.Email)
	assert.NoError(t, err)
	assert.Equal(t, createdProfile, foundProfile)
}
