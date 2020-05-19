package models

type User struct {
	ID             int
	Email          string
	HashedPassword string
	Active         bool
}
