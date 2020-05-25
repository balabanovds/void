package domain

import "errors"

var (
	ErrNotFound           = errors.New("entry not found")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrAlreadyExists      = errors.New("already exists")
	ErrDependencyNotFound = errors.New("dependent entries not found")
)
