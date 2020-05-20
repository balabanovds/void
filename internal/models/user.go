package models

import "time"

// User ...
type User struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	HashedPassword []byte    `json:"-"`
	Active         bool      `json:"active"`
	Created        time.Time `json:"created"`
}
