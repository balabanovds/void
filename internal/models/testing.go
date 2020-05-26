package models

import (
	"testing"
)

// TestUser is a helper function just for testing
func TestUser(t *testing.T) *User {
	t.Helper()

	return &User{
		Email:          "test@mail.com",
		Password:       "password123",
		HashedPassword: []byte("passw"),
	}
}

func TestProfile(t *testing.T) Profile {
	t.Helper()

	return Profile{
		Email:     "test@mail.com",
		FirstName: "Vasya",
		LastName:  "Pupkin",
		Phone:     "79212223344",
		Position:  "tech lead",
		CompanyID: 123456,
		ZCode:     "ZA1234",
		Ru: ProfileRu{
			FirstName: "Вася",
		},
	}
}

func TestNewProfile(t *testing.T) NewProfile {
	t.Helper()

	return NewProfile{
		Email:     "",
		FirstName: "Vasya",
		LastName:  "Pupkin",
		Phone:     "79212223344",
		Position:  "tech lead",
		CompanyID: 123456,
		ZCode:     "ZA1234",
		Ru: ProfileRu{
			FirstName:  "Вася",
			Patronymic: "Петрович",
			LastName:   "Пупкин",
			Position:   "Тех лид",
		},
	}
}
