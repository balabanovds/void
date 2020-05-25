package models

import (
	"database/sql"
	"time"
)

type Profile struct {
	ID           int            `json:"id"`
	FirstName    string         `json:"first_name_en"`
	LastName     string         `json:"last_name_en"`
	Email        string         `json:"email"`
	Phone        string         `json:"phone"`
	Position     string         `json:"position_en"`
	CompanyID    int            `json:"company_id"`
	ZCode        string         `json:"z_code"`
	ManagerEmail sql.NullString `json:"manager_email"`
	Role         Role           `json:"role"`
	Ru           ProfileRu      `json:"ru"`
	ModifiedAt   time.Time      `json:"modified_at_en"`
}

// NewProfile describes which fields we expect to get from client
type NewProfile struct {
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	FirstName string    `json:"first_name_en"`
	LastName  string    `json:"last_name_en"`
	Position  string    `json:"position_en"`
	CompanyID int       `json:"company_id"`
	ZCode     string    `json:"z_code"`
	Ru        ProfileRu `json:"ru"`
}

// UpdateProfile describes which fields we can update
type UpdateProfile struct {
	FirstName    string    `json:"first_name_en"`
	LastName     string    `json:"last_name_en"`
	Phone        string    `json:"phone"`
	Position     string    `json:"position_en"`
	CompanyID    int       `json:"company_id"`
	ZCode        string    `json:"z_code"`
	ManagerEmail string    `json:"manager_email"`
	Role         Role      `json:"role"`
	Ru           ProfileRu `json:"ru"`
}

type ProfileRu struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Patronymic string `json:"patronymic"`
	Position   string `json:"position"`
}

type Role struct {
	ID    int
	Value string
}

var (
	Engineer = Role{
		ID:    1,
		Value: "engineer",
	}
	Manager = Role{
		ID:    2,
		Value: "manager",
	}
	Admin = Role{
		ID:    99,
		Value: "admin",
	}
)
