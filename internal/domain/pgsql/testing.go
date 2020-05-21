package pgsql

import (
	"fmt"
	"github.com/balabanovds/void/internal/models"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

type TestSuite struct {
	Storage *Storage
	clear   func(...string)
	User    *models.User
}

func NewTestSuite(t *testing.T) TestSuite {
	t.Helper()

	s, cl := TestDB(t)
	return TestSuite{
		Storage: s,
		clear:   cl,
		User:    models.TestUser(t),
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
		_, err := s.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", strings.Join(tables, ", ")))
		if err != nil {
			t.Fatal(err)
		}
	}

}

func getURL(t *testing.T) string {
	t.Helper()

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "host=localhost port=5432 user=void password=void123 dbname=void_test sslmode=disable"
	}
	return url
}
