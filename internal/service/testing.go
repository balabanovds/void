package service

import (
	"github.com/balabanovds/void/internal/domain/pgsql"
	"github.com/balabanovds/void/internal/models"
	"github.com/rs/zerolog"
	"testing"
)

type TestSuite struct {
	Service *Service
	Storage pgsql.TestSuite
}

func NewTestSuite(t *testing.T) TestSuite {
	t.Helper()

	ts := pgsql.NewTestSuite(t)
	return TestSuite{
		Service: New(ts.Storage, zerolog.Nop()),
		Storage: ts,
	}
}

func (s *TestSuite) CreateUser(email, password string) (models.User, error) {
	return s.Service.Users().Create(email, password, password)
}

func (s *TestSuite) Close() {
	s.Storage.Close()
}
