package service

import (
	"github.com/balabanovds/void/internal/domain"
	"github.com/rs/zerolog"
)

// CtxKey service context keys
type CtxKey uint8

const (
	// CtxIsAdmin bool
	CtxIsAdmin CtxKey = iota
	// CtxEmail string - email of currently logged in user
	CtxEmail
)

// Service ...
type Service struct {
	storage     domain.Storage
	userservice *UserService
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
	if s.userservice == nil {
		s.userservice = newUserService(s)
	}
	return s.userservice
}
