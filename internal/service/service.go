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
}

// New service
func New(storage domain.Storage, logger zerolog.Logger) *Service {
	return &Service{
		storage: storage,
		log:     logger.With().Str("service", "SERVICE").Logger(),
	}
}

// Users service
func (s *Service) Users() *UserService {
	if s.userService == nil {
		s.userService = newUserService(s)
	}
	return s.userService
}
