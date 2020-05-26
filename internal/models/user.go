package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"time"
)

// User ...
type User struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	HashedPassword []byte    `json:"-"`
	Active         bool      `json:"active"`
	Created        time.Time `json:"created"`
}

// NewUser explains what we expect to receive in domain while creating new user
type NewUser struct {
	Email          string `json:"email"`
	HashedPassword []byte `json:"-"`
}

// UpdateUser expected fields for update in domain
type UpdatedUser struct {
	HashedPassword []byte `json:"-"`
	Active         bool   `json:"active"`
}

func (u *User) Sanitize() {
	u.Password = ""
}

func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(6, 100)),
	)
}
