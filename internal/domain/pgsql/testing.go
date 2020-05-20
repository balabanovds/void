package pgsql

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

// TestDB helper returning storage and cleanup function that receive table names to truncate
func TestDB(t *testing.T, url string) (*Storage, func(...string)) {
	t.Helper()

	s := New(nil, zerolog.Nop())
	err := s.openURL(url)
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
