package sqlrepo

import (
	"database/sql"

	"github.com/balabanovds/void/internal/domain"
	"github.com/rs/zerolog"
)

// Storage main struct containing all SQL repositories implementations
type Storage struct {
	db       *sql.DB
	userRepo *userRepo
	logger   *zerolog.Logger
}

// New Storage
func New(db *sql.DB, logger *zerolog.Logger) *Storage {
	return &Storage{db: db, logger: logger}
}

// Users return users repository SQL implementation
func (s *Storage) Users() domain.UserRepo {
	if s.userRepo == nil {
		s.userRepo = newUserRepo(s.db)
	}
	return s.userRepo
}
