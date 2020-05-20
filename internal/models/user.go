package models

import "time"

// User ...
type User struct {
	ID             int
	Email          string
	HashedPassword string
	Active         bool
	Created        time.Time
}
