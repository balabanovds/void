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
	repo  domain.UserRepo
	log   zerolog.Logger
	debug zerolog.Logger
}

func newUserService(service *Service) *UserService {
	return &UserService{
		repo:  service.storage.Users(),
		log:   service.log,
		debug: service.debug,
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
		s.debug.Error().Msg(err.Error())
		return models.User{}, err
	}

	return s.repo.Create(email, hash)
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
	return s.repo.Get(email)
}

// IsAdmin checks user rights
func (s *UserService) IsAdmin(email string) (bool, error) {
	//TODO when profiles will be ready
	return false, nil
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
		s.debug.Error().Msg(err.Error())
		return err
	}

	u, err := s.repo.Get(email)
	if err != nil {
		return err
	}

	return s.repo.Update(&u, hash, u.Active)
}

// ToggleActive state (only admin allowed)
// Client errors:
// service.ErrNotAllowed
func (s *UserService) ToggleActive(ctx context.Context, email string) error {
	if !ctxHelper.IsAdmin(ctx) {
		return ErrNotAllowed
	}

	u, err := s.repo.Get(email)
	if err != nil {
		return err
	}

	err = s.repo.Update(&u, nil, !u.Active)
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

	s.repo.Delete(email)
	return nil
}
