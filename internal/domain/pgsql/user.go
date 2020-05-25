package pgsql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/balabanovds/void/internal/domain"
	"github.com/balabanovds/void/internal/models"
)

type userRepo struct {
	//db *sql.DB
	st *Storage
}

func newUserRepo(storage *Storage) *userRepo {
	return &userRepo{st: storage}
}

func (r *userRepo) Create(newUser models.NewUser) (models.User, error) {
	u := models.User{
		Email:          newUser.Email,
		HashedPassword: newUser.HashedPassword,
	}

	if err := r.st.db.QueryRow(
		"insert into users (email, hashed_password, created) "+
			"values ($1, $2, now()) returning id, active, created",
		newUser.Email, newUser.HashedPassword,
	).Scan(&u.ID, &u.Active, &u.Created); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return models.User{}, domain.ErrDuplicateEmail
		}
		return models.User{}, err
	}

	return u, nil
}

func (r *userRepo) Get(email string) (models.User, error) {
	u := models.User{}
	if err := r.st.db.QueryRow(
		"select id, email, hashed_password, created, active from users where email = $1",
		email,
	).Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Created, &u.Active); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, domain.ErrNotFound
		}
		return models.User{}, err
	}
	return u, nil
}

func (r *userRepo) Update(user *models.User, newHashPassword []byte, newActive bool) error {
	if err := r.st.db.QueryRow(
		"UPDATE users SET active = $2, hashed_password = $3 "+
			"where id = $1 returning active, hashed_password",
		user.ID, newActive, newHashPassword,
	).Scan(&user.Active, &user.HashedPassword); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) Delete(email string) {
	r.st.db.QueryRow("DELETE FROM users WHERE email = $1", email)
}
