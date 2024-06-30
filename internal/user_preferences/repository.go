package userpreferences

import (
	"database/sql"
)

type UserPreferencesRepository struct {
	db *sql.DB
}

func NewUserPreferencesRepository(db *sql.DB) *UserPreferencesRepository {
	return &UserPreferencesRepository{db: db}
}

func (r *UserPreferencesRepository) CreateTable() error {
	return nil
}

func (r *UserPreferencesRepository) DropTable() error {
	return nil
}
