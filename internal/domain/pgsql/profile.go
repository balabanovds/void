package pgsql

import (
	"database/sql"
	"errors"
	"github.com/balabanovds/void/internal/domain"
	"github.com/balabanovds/void/internal/models"
	"strings"
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

func (r *profileRepo) Get(email string) (models.Profile, error) {
	p := models.Profile{}

	if err := r.db.QueryRow(
		"SELECT p.id, p.email, p.first_name, p.last_name, p.position, p.phone, p.company_id, "+
			"p.z_code, p.manager_email, r.id, r.value, p.modified_at, "+
			"ru.first_name, ru.last_name, ru.patronymic, ru.position "+
			"FROM profiles AS p " +
			"JOIN roles AS r ON p.role_id = r.id " +
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

func (r *profileRepo) Update(profile *models.Profile, upd models.UpdatedProfile) error {
	return nil
}