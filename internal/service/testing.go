package service

import (
	"testing"

	"github.com/balabanovds/void/internal/domain/pgsql"
	"github.com/balabanovds/void/internal/models"
	"github.com/rs/zerolog"
)

// TestSuite for service
type TestSuite struct {
	Service *Service
	Storage pgsql.TestSuite
}

// NewTestSuite ...
func NewTestSuite(t *testing.T) TestSuite {
	t.Helper()

	ts := pgsql.NewTestSuite(t)
	return TestSuite{
		Service: New(ts.Storage, zerolog.Nop()),
		Storage: ts,
	}
}

// CreateDefaultUser creates test user
func (s *TestSuite) CreateUser(email, password string) (models.User, error) {
	return s.Service.Users().Create(email, password, password)
}

// Close test suite
func (s *TestSuite) Close() {
	s.Storage.Close()
}
