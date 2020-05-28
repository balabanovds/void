package pgsql

import (
	"testing"

	"github.com/balabanovds/void/internal/domain"
	"github.com/balabanovds/void/internal/models"
	"github.com/stretchr/testify/assert"
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

	profile := ts.CreateDefaultProfile(t)

	_, err := ts.Storage.Profiles().Get(profile.Email + "wrong")
	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrNotFound.Error())

	foundProfile, err := ts.Storage.Profiles().Get(profile.Email)
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

	// check delete manager email
	upd.ManagerEmail = ""

	err = ts.Storage.Profiles().Update(&p, upd)
	assert.NoError(t, err)
	assert.False(t, p.ManagerEmail.Valid)
	mgr, err := p.ManagerEmail.Value()
	assert.NoError(t, err)
	assert.Empty(t, mgr)

	// check wrong manager email
	upd.ManagerEmail = manProfile.Email + "wrong"
	err = ts.Storage.Profiles().Update(&p, upd)
	assert.Error(t, err)
	assert.EqualError(t, err, domain.ErrDependencyNotFound.Error())
}

func TestProfileRepo_GetAll(t *testing.T) {
	ts := NewTestSuite(t)
	defer ts.Close()

	ps := ts.Storage.Profiles().GetAll()
	assert.Nil(t, ps)

	_ = ts.CreateProfile(t, "mail1")
	_ = ts.CreateProfile(t, "mail2")

	ps = ts.Storage.Profiles().GetAll()
	assert.Len(t, ps, 2)
}
