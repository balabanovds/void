package pgsql

import (
	"database/sql"
	"fmt"

	"github.com/balabanovds/void/internal/domain"
	"github.com/rs/zerolog"

	_ "github.com/lib/pq" // ...
)

// Storage main struct containing all SQL repositories implementations
type Storage struct {
	config      *domain.Config
	db          *sql.DB
	userRepo    *userRepo
	profileRepo *profileRepo
	log         zerolog.Logger
}

// New Storage
func New(config *domain.Config, logger zerolog.Logger) *Storage {
	return &Storage{
		config: config,
		log:    logger.With().Str("storage", "psql").Logger(),
	}
}

// Users return users repository SQL implementation
func (s *Storage) Users() domain.UserRepo {
	if s.db == nil {
		s.log.Fatal().Msg("sql connection not opened")
	}
	if s.userRepo == nil {
		s.userRepo = newUserRepo(s)
	}
	return s.userRepo
}

// Profiles provides access to Profiles repository
func (s *Storage) Profiles() domain.ProfileRepo {
	if s.db == nil {
		s.log.Fatal().Msg("sql connection not opened")
	}
	if s.profileRepo == nil {
		s.profileRepo = newProfileRepo(s.db)
	}
	return s.profileRepo
}

func (s *Storage) openURL(url string) error {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}
	s.db = db

	return nil

}

// Open Postgres SQL pool
func (s *Storage) Open() error {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.config.SQL.Host,
		s.config.SQL.Port,
		s.config.SQL.User,
		s.config.SQL.Password,
		s.config.SQL.DBName)

	if err := s.openURL(dsn); err != nil {
		return err
	}

	s.log.Info().Msgf("connected to '%s' at %s:%d",
		s.config.SQL.DBName,
		s.config.SQL.Host,
		s.config.SQL.Port,
	)

	return nil
}

// Close SQL pool
func (s *Storage) Close() {
	if s.db != nil {
		_ = s.db.Close()
	}
}
