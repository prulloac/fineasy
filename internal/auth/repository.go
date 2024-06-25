package auth

import (
	"database/sql"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (a *AuthRepository) InsertUser(u *User) error {
	// check if user already exists
	var id int
	err := a.db.QueryRow(`
	SELECT 
		id
	FROM users
	WHERE email = $1
	`, u.Email).Scan(&id)

	if err == sql.ErrNoRows {
		_, err = a.db.Exec(`
		INSERT INTO users (hash, username, email)
		VALUES ($1, $2, $3, $4, $5, $6)
		`, u.Hash, u.Username, u.Email)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}
