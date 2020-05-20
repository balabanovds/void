package domain

import "github.com/balabanovds/void/internal/models"

// Storage ...
type Storage interface {
	Users() UserRepo
	Open() error
	Close()
}

// UserRepo is a contract for user reposotories implementations
type UserRepo interface {
	// Create user and returning created user from DB
	Create(email string, hashPassword []byte) (models.User, error)
	// GetByEmail ...
	Get(email string) (models.User, error)
	// Update user HashedPassword and Active state only
	Update(user *models.User, newHashPassword []byte, active bool) error
	// Delete by id
	Delete(email string)
}
