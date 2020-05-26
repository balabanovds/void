package service

import (
	"context"

	"github.com/balabanovds/void/internal/domain"
	"github.com/balabanovds/void/internal/models"
	"github.com/balabanovds/void/internal/server/ctxHelper"
)

// CreateProfile can be executed only by self
// Client errors:
// service.ErrNotAllowed,
// domain.ErrAlreadyExists,
// domain.ErrDependencyNotFound
func (s *UserService) CreateProfile(ctx context.Context, newProfile models.NewProfile) (models.Profile, error) {
	// check if user self creating profile
	if !ctxHelper.IsEmailMatch(ctx, newProfile.Email) {
		return models.Profile{}, ErrNotAllowed
	}

	// check if user exists
	_, err := s.GetByEmail(newProfile.Email)
	if err != nil {
		return models.Profile{}, domain.ErrDependencyNotFound
	}

	return s.storage.Profiles().Create(newProfile)
}

// GetProfile ...
// Client errors:
// domain.ErrNotFound,
func (s *UserService) GetProfile(email string) (models.Profile, error) {
	profile, err := s.storage.Profiles().Get(email)
	if err != nil {
		return models.Profile{}, err
	}
	return profile, nil
}

// GetAllProfiles ...
// Client errors:
// domain.ErrNotFound
func (s *UserService) GetAllProfiles() []models.Profile {
	return s.storage.Profiles().GetAll()
}

// UpdateProfile only by self or admin
// Client errors:
// service.ErrNotAllowed
// domain.ErrNotFound
// domain.ErrDependencyNotFound
func (s *UserService) UpdateProfile(ctx context.Context, p *models.Profile, upd models.UpdateProfile) error {
	// check if user self or admin
	if !ctxHelper.IsEmailMatch(ctx, p.Email) && !ctxHelper.IsAdmin(ctx) {
		return ErrNotAllowed
	}

	// check if profile exists in db
	_, err := s.storage.Profiles().Get(p.Email)
	if err != nil {
		return err
	}

	if err := s.storage.Profiles().Update(p, upd); err != nil {
		return err
	}

	return nil
}

// ChangeRole only for admin
// Client errors:
// service.ErrNotAllowed
// domain.ErrNotFound
// domain.ErrDependencyNotFound
func (s *UserService) ChangeRole(ctx context.Context, p *models.Profile, role models.Role) error {
	if !ctxHelper.IsAdmin(ctx) {
		return ErrNotAllowed
	}

	// check if role exists in db first
	_, err := s.storage.Profiles().Get(p.Email)
	if err != nil {
		return err
	}

	err = s.storage.Profiles().Update(p, models.UpdateProfile{Role: role})
	if err != nil {
		return err
	}
	return nil
}
