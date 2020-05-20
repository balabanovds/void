package service

import "errors"

var (
	// ErrFailedAuthenticate ...
	ErrFailedAuthenticate = errors.New("failed to authenticate")
	// ErrEmailContext ...
	ErrEmailContext = errors.New("context does not contain 'email'")
	// ErrNotAllowed ...
	ErrNotAllowed = errors.New("action not allowed")
	// ErrPasswdNotMatch ...
	ErrPasswdNotMatch = errors.New("passwords not match")
)
