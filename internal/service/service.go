package service

import (
	"github.com/balabanovds/void/internal/domain"
	"github.com/rs/zerolog"
)

// Service ...
type Service struct {
	storage     domain.Storage
	userService *UserService
	log         zerolog.Logger
	debug       zerolog.Logger
}

// New service
func New(storage domain.Storage, logger zerolog.Logger) *Service {
	l := logger.With().Str("service", "SERVICE").Logger()
	return &Service{
		storage: storage,
		log:     l,
		debug:   l.With().Caller().Stack().Logger(),
	}
}

// Users service
func (s *Service) Users() *UserService {
	if s.userService == nil {
		s.userService = newUserService(s)
	}
	return s.userService
}
