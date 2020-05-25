package pgsql

import (
	"fmt"
	"github.com/balabanovds/void/internal/models"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

type TestSuite struct {
	Storage *Storage
	clear   func(...string)
	User    *models.User
	Profile models.Profile
}

func NewTestSuite(t *testing.T) TestSuite {
	t.Helper()

	s, cl := TestDB(t)
	return TestSuite{
		Storage: s,
		clear:   cl,
		User:    models.TestUser(t),
		Profile: models.TestProfile(t),
	}
}

func (ts *TestSuite) Close() {
	ts.clear("users")
	ts.Storage.Close()
}

// TestDB helper returning storage and cleanup function that receive table names to truncate
func TestDB(t *testing.T) (*Storage, func(...string)) {
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
		url = "host=balabanov.sknt.ru port=5432 user=void password=@ws3ed4rf dbname=void_test sslmode=disable"
	}
	return url
}
