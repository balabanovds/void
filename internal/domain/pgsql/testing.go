package pgsql

import (
	"fmt"
	"os"
	"testing"

	"github.com/balabanovds/void/internal/models"

	"github.com/rs/zerolog"
)

// TestSuite ...
type TestSuite struct {
	Storage *Storage
	clear   func(...string)
	User    *models.User
	Profile models.Profile
}

// NewTestSuite helper
func NewTestSuite(t *testing.T) TestSuite {
	t.Helper()

	s, cl := testDB(t)
	return TestSuite{
		Storage: s,
		clear:   cl,
		User:    models.TestUser(t),
		Profile: models.TestProfile(t),
	}
}

// Close DB but before clear "users" table
func (ts *TestSuite) Close() {
	ts.clear("users")
	ts.Storage.Close()
}

func (ts *TestSuite) CreateUser(t *testing.T, email string) models.User {
	t.Helper()
	if email == "" {
		email = ts.User.Email
	}

	newUser := models.NewUser{
		Email:          email,
		HashedPassword: ts.User.HashedPassword,
	}

	u, err := ts.Storage.Users().Create(newUser)
	if err != nil {
		t.Fatal(err)
	}
	return u
}

func (ts *TestSuite) CreateDefaultUser(t *testing.T) models.User {
	return ts.CreateUser(t, "")
}

func (ts *TestSuite) CreateDefaultProfile(t *testing.T) models.Profile {
	return ts.CreateProfile(t, "")
}

func (ts *TestSuite) CreateProfile(t *testing.T, email string) models.Profile {
	if email == "" {
		email = ts.User.Email
	}

	_ = ts.CreateUser(t, email)

	newProfile := models.TestNewProfile(t)
	newProfile.Email = email
	createdProfile, err := ts.Storage.Profiles().Create(newProfile)
	if err != nil {
		t.Fatal(err)
	}
	return createdProfile
}

// testDB helper returning storage and cleanup function that receive table names to truncate
func testDB(t *testing.T) (*Storage, func(...string)) {
	t.Helper()

	s := New(nil, zerolog.Nop())
	err := s.openURL(getURL(t))
	if err != nil {
		t.Fatal(err)
	}

	return s, func(tables ...string) {
		defer s.Close()
		for _, table := range tables {
			_, err := s.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
			if err != nil {
				t.Fatal(err)
			}
		}
	}

}

func getURL(t *testing.T) string {
	t.Helper()

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "host=localhost port=5432 user=void password=void123 dbname=void_test sslmode=disable"
		//url = "host=balabanov.sknt.ru port=5432 user=void password=@ws3ed4rf dbname=void_test sslmode=disable"
	}
	return url
}
