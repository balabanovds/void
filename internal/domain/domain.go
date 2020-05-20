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
	Create(email, password string) (models.User, error)
	// Get by id
	Get(id int) (models.User, error)
	// GetByEmail ...
	GetByEmail(email string) (models.User, error)
	// Update user HashedPassword and Active state only
	Update(user *models.User, newHashPassword string, active bool) error
	// Delete by id
	Delete(id int)
}
