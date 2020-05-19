package sqlrepo

import "database/sql"

type userRepo struct {
	db *sql.DB
}

func newUserRepo(db *sql.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(email, password string) error {
	return nil
}
