package pgsql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/balabanovds/void/internal/domain"
	"github.com/balabanovds/void/internal/models"
)

type profileRepo struct {
	db *sql.DB
}

func newProfileRepo(db *sql.DB) *profileRepo {
	return &profileRepo{db: db}
}

// Create new profile en and ru
func (r *profileRepo) Create(pr models.NewProfile) (models.Profile, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return models.Profile{}, err
	}

	profile := models.Profile{
		Email:     pr.Email,
		FirstName: pr.FirstName,
		LastName:  pr.LastName,
		Position:  pr.Position,
		Phone:     pr.Phone,
		ZCode:     pr.ZCode,
		CompanyID: pr.CompanyID,
		Role:      models.Engineer,
		Ru: models.ProfileRu{
			FirstName:  pr.Ru.FirstName,
			LastName:   pr.Ru.LastName,
			Patronymic: pr.Ru.Patronymic,
			Position:   pr.Ru.Position,
		},
	}

	if err := tx.QueryRow("INSERT INTO profiles "+
		"(email, first_name, last_name, phone, position, company_id, z_code, role_id, modified_at) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now()) RETURNING id, modified_at",
		profile.Email,
		profile.FirstName,
		profile.LastName,
		profile.Phone,
		profile.Position,
		profile.CompanyID,
		profile.ZCode,
		profile.Role.ID).
		Scan(&profile.ID, &profile.ModifiedAt); err != nil {

		_ = tx.Rollback()
		if strings.Contains(err.Error(), "unique") {
			return models.Profile{}, domain.ErrAlreadyExists
		}
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return models.Profile{}, domain.ErrDependencyNotFound
		}
		return models.Profile{}, err
	}

	if _, err := tx.Exec(
		"INSERT INTO profiles_ru (profile_id, first_name, last_name, patronymic, position) "+
			"VALUES ($1, $2, $3, $4, $5)",
		profile.ID,
		profile.Ru.FirstName,
		profile.Ru.LastName,
		profile.Ru.Patronymic,
		profile.Ru.Position); err != nil {
		if strings.Contains(err.Error(), "unique") {
			return models.Profile{}, domain.ErrAlreadyExists
		}
		return models.Profile{}, err
	}

	err = tx.Commit()
	if err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}

// Get ...
func (r *profileRepo) Get(email string) (models.Profile, error) {
	p := models.Profile{}

	if err := r.db.QueryRow(
		"SELECT p.id, p.email, p.first_name, p.last_name, p.position, p.phone, p.company_id, "+
			"p.z_code, p.manager_email, r.id, r.value, p.modified_at, "+
			"ru.first_name, ru.last_name, ru.patronymic, ru.position "+
			"FROM profiles AS p "+
			"JOIN roles AS r ON p.role_id = r.id "+
			"JOIN profiles_ru AS ru ON ru.profile_id = p.id "+
			"WHERE p.email = $1", email).
		Scan(
			&p.ID,
			&p.Email,
			&p.FirstName,
			&p.LastName,
			&p.Position,
			&p.Phone,
			&p.CompanyID,
			&p.ZCode,
			&p.ManagerEmail,
			&p.Role.ID,
			&p.Role.Value,
			&p.ModifiedAt,
			&p.Ru.FirstName,
			&p.Ru.LastName,
			&p.Ru.Patronymic,
			&p.Ru.Position,
		); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Profile{}, domain.ErrNotFound
		}
		return models.Profile{}, err
	}

	return p, nil
}

// Update profile that shared down the stack with updated profile
// Client errors
func (r *profileRepo) Update(profile *models.Profile, upd models.UpdateProfile) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if err := tx.QueryRow(
		"UPDATE profiles SET first_name = $2, last_name = $3, position = $4, "+
			"phone = $5, company_id = $6, z_code = $7, manager_email = $8, "+
			"role_id = $9, modified_at = now() WHERE email = $1 RETURNING modified_at",
		profile.Email,
		upd.FirstName,
		upd.LastName,
		upd.Position,
		upd.Phone,
		upd.CompanyID,
		upd.ZCode,
		upd.ManagerEmail,
		upd.Role.ID).Scan(&profile.ModifiedAt); err != nil {
		_ = tx.Rollback()
		return err
	}

	if _, err := tx.Exec(
		"UPDATE profiles_ru SET first_name = $2, last_name = $3, patronymic = $4, position = $5 "+
			"WHERE profile_id = $1",
		profile.ID,
		upd.Ru.FirstName,
		upd.Ru.LastName,
		upd.Ru.Patronymic,
		upd.Ru.Position); err != nil {
		_ = tx.Rollback()
		return err
	}

	// if db updated without errors than update profile
	profile.FirstName = upd.FirstName
	profile.LastName = upd.LastName
	profile.Position = upd.Position
	profile.Phone = upd.Phone
	profile.CompanyID = upd.CompanyID
	profile.ZCode = upd.ZCode
	profile.ManagerEmail = sql.NullString{
		String: upd.ManagerEmail,
		Valid:  true,
	}
	profile.Role = upd.Role
	profile.Ru = upd.Ru

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
