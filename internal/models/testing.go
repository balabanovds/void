package models

import "testing"

// TestUser is a helper function just for testing
func TestUser(t *testing.T) *User {
	t.Helper()

	return &User{
		Email:          "test@mail.com",
		Password:       "password123",
		HashedPassword: []byte("passw"),
	}
}
