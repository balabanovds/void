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

	user := ts.CreateDefaultUser(t)

	newProfile := models.TestNewProfile(t)

	newProfile.Email = user.Email + "wrong"
	_, err := ts.Storage.Profiles().Create(newProfile)
	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrDependencyNotFound.Error())

	newProfile.Email = user.Email

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

	user := ts.CreateDefaultUser(t)

	profile := ts.CreateDefaultProfile(t)

	_, err := ts.Storage.Profiles().Get(user.Email + "wrong")
	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrNotFound.Error())

	foundProfile, err := ts.Storage.Profiles().Get(user.Email)
	assert.NoError(t, err)
	assert.Equal(t, profile, foundProfile)
}

func TestProfileRepo_Update(t *testing.T) {
	ts := NewTestSuite(t)
	defer ts.Close()

	p := ts.CreateDefaultProfile(t)
	manProfile := ts.CreateProfile(t, "manager@mail.org")

	upd := models.UpdateProfile{
		FirstName:    p.FirstName + "upd",
		LastName:     p.LastName + "upd",
		Phone:        p.Phone + "upd",
		Position:     p.Position + "upd",
		CompanyID:    p.CompanyID + 1,
		ZCode:        p.ZCode + "upd",
		ManagerEmail: manProfile.Email,
		Role: models.Role{
			ID:    99,
			Value: "admin",
		},
		Ru: p.Ru,
	}

	initProfile := p

	err := ts.Storage.Profiles().Update(&p, upd)
	assert.NoError(t, err)
	assert.NotEqual(t, initProfile, p)

}
