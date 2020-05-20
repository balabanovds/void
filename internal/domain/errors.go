package domain

import "errors"

var (
	// ErrNotFound ...
	ErrNotFound = errors.New("entry not found")
	// ErrDuplicateEmail ...
	ErrDuplicateEmail = errors.New("duplicate email")
	// ErrInvalidCredentials ...
	ErrInvalidCredentials = errors.New("invalid credentials")
)
