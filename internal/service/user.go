package service

import (
	"context"

	"github.com/balabanovds/void/internal/server/ctxHelper"

	"github.com/balabanovds/void/internal/domain"
	"github.com/balabanovds/void/internal/models"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

// UserService ...
type UserService struct {
	storage domain.Storage
	log     zerolog.Logger
	debug   func(err error)
}

func newUserService(service *Service) *UserService {
	return &UserService{
		storage: service.storage,
		log:     service.log,
		debug:   service.debugLog,
	}
}

// Create new user.
// Client errors:
// service.ErrPasswdNotMatch,
// domain.ErrDuplicateEmail
func (s *UserService) Create(email, password, confirmPassword string) (models.User, error) {
	if password != confirmPassword {
		return models.User{}, ErrPasswdNotMatch
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.debug(err)
		return models.User{}, err
	}

	return s.storage.Users().Create(models.NewUser{Email: email, HashedPassword: hash})
}

// Authenticate while logging user
// Client errors:
// domain.ErrNotFound,
// service.ErrFailedAuthenticate
func (s *UserService) Authenticate(email, password string) (models.User, error) {
	u, err := s.GetByEmail(email)
	if err != nil {
		return models.User{}, err
	}

	if !u.Active {
		return models.User{}, ErrFailedAuthenticate
	}

	err = bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password))
	if err != nil {
		return models.User{}, ErrFailedAuthenticate
	}

	return u, nil
}

// GetByEmail ...
// Client errors:
// domain.ErrNotFound
func (s *UserService) GetByEmail(email string) (models.User, error) {
	return s.storage.Users().Get(email)
}

// UpdatePassword only self can do
// Client errors:
// service.ErrNotAllowed
func (s *UserService) UpdatePassword(ctx context.Context, email, password string) error {
	if !ctxHelper.IsEmailMatch(ctx, email) {
		return ErrNotAllowed
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.debug(err)
		return err
	}

	u, err := s.GetByEmail(email)
	if err != nil {
		return err
	}

	return s.storage.Users().Update(&u, hash, u.Active)
}

// ToggleActive state (only admin allowed)
// Client errors:
// service.ErrNotAllowed
func (s *UserService) ToggleActive(ctx context.Context, email string) error {
	if !ctxHelper.IsAdmin(ctx) {
		return ErrNotAllowed
	}

	u, err := s.GetByEmail(email)
	if err != nil {
		return err
	}

	err = s.storage.Users().Update(&u, nil, !u.Active)
	if err != nil {
		return err
	}
	return nil
}

// Delete can do self or admin
// Client errors:
// service.ErrNotAllowed
func (s *UserService) Delete(ctx context.Context, email string) error {
	if !ctxHelper.IsEmailMatch(ctx, email) && !ctxHelper.IsAdmin(ctx) {
		return ErrNotAllowed
	}

	s.storage.Users().Delete(email)
	return nil
}
