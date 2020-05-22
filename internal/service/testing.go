package service

import (
	"github.com/balabanovds/void/internal/domain/pgsql"
	"github.com/balabanovds/void/internal/models"
	"github.com/rs/zerolog"
	"testing"
)

type ServiceTestSuite struct {
	Service *Service
	storage pgsql.TestSuite
}

func NewTestSuite(t *testing.T) ServiceTestSuite {
	t.Helper()

	ts := pgsql.NewTestSuite(t)
	return ServiceTestSuite{
		Service: New(ts.Storage, zerolog.Nop()),
		storage: ts,
	}
}

func (s *ServiceTestSuite) CreateUser(email, password string) (models.User, error) {
	return s.Service.Users().Create(email, password, password)
}

func (s *ServiceTestSuite) Close() {
	s.storage.Close()
}
