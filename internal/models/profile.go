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

func (u UpdateProfile) CopyToProfile(profile *Profile) {
	if u.FirstName != "" {
		profile.FirstName = u.FirstName
	}
	if u.LastName != "" {
		profile.LastName = u.LastName
	}
	if u.Position != "" {
		profile.Position = u.Position
	}
	if u.Phone != "" {
		profile.Phone = u.Phone
	}
	if u.CompanyID != 0 {
		profile.CompanyID = u.CompanyID
	}
	if u.ZCode != "" {
		profile.ZCode = u.ZCode
	}
	profile.ManagerEmail = sql.NullString{
		String: u.ManagerEmail,
		Valid:  u.ManagerEmail != "",
	}
	if u.Role.ID != 0 && u.Role.Value != "" {
		profile.Role = u.Role
	}
	if u.Ru.FirstName != "" {
		profile.Ru.FirstName = u.Ru.FirstName
	}
	if u.Ru.LastName != "" {
		profile.Ru.LastName = u.Ru.LastName
	}
	if u.Ru.Patronymic != "" {
		profile.Ru.Patronymic = u.Ru.Patronymic
	}
	if u.Ru.Position != "" {
		profile.Ru.Position = u.Ru.Position
	}
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
