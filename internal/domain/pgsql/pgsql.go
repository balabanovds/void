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
	config   *domain.Config
	db       *sql.DB
	userRepo *userRepo
	log      zerolog.Logger
	debug    zerolog.Logger
}

// New Storage
func New(config *domain.Config, logger zerolog.Logger) *Storage {
	l := logger.With().Str("storage", "PSQL").Logger()
	return &Storage{
		config: config,
		log:    l,
		debug:  l.With().Caller().Stack().Logger(),
	}
}

// Users return users repository SQL implementation
func (s *Storage) Users() domain.UserRepo {
	if s.db == nil {
		s.debug.Fatal().Msg("sql connection not opened")
	}
	if s.userRepo == nil {
		s.userRepo = newUserRepo(s.db)
	}
	return s.userRepo
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

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}
	s.db = db

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
		s.db.Close()
	}
}
