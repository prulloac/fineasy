package auth

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (a *AuthRepository) CreateTable() error {
	_, err := a.DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		hash uuid NOT NULL,
		username VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		validated_at TIMESTAMP,
		disabled BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS internal_logins (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL references users(id) ON DELETE CASCADE,
		email VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		password_salt VARCHAR(255) NOT NULL,
		algorithm INTEGER NOT NULL,
		password_last_updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		login_attempts INT NOT NULL DEFAULT 0,
		last_login_attempt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		last_login_success TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS login_tokens (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL references users(id) ON DELETE CASCADE,
		token VARCHAR(255) NOT NULL,
		token_type INTEGER NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		used_at TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS external_login_providers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		type INTEGER NOT NULL,
		endpoint VARCHAR(255) NOT NULL,
		enabled BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS external_logins (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL references users(id) ON DELETE CASCADE,
		provider_id INT NOT NULL references external_login_providers(id) ON DELETE CASCADE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS external_login_tokens (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL references users(id) ON DELETE CASCADE,
		provider_id INT NOT NULL references external_login_providers(id) ON DELETE CASCADE,
		login_ip VARCHAR(255) NOT NULL,
		user_agent VARCHAR(255) NOT NULL,
		logged_in_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		token TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS user_sessions (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL references users(id) ON DELETE CASCADE,
		login_ip VARCHAR(255) NOT NULL,
		user_agent VARCHAR(255) NOT NULL,
		logged_in_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		logged_out_at TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE UNIQUE INDEX IF NOT EXISTS users_email_idx ON users (email);
	CREATE UNIQUE INDEX IF NOT EXISTS internal_logins_email_idx ON internal_logins (email);
	CREATE UNIQUE INDEX IF NOT EXISTS external_login_providers_name_idx ON external_login_providers (name, type, endpoint);
	CREATE UNIQUE INDEX IF NOT EXISTS external_logins_user_id_provider_id_idx ON external_logins (user_id, provider_id);
	`)
	return err
}

func (a *AuthRepository) DropTable() error {
	_, err := a.DB.Exec(`
	DROP TABLE IF EXISTS user_sessions;
	DROP TABLE IF EXISTS external_login_tokens;
	DROP TABLE IF EXISTS external_logins;
	DROP TABLE IF EXISTS external_login_providers;
	DROP TABLE IF EXISTS login_tokens;
	DROP TABLE IF EXISTS internal_logins;
	DROP TABLE IF EXISTS users;
	`)
	return err
}

func (a *AuthRepository) getUserID(email string) (int, error) {
	var uid int
	err := a.DB.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&uid)
	return uid, err
}

func (a *AuthRepository) getSaltAndAlgorithmForUser(uid int) (string, Algorithm, error) {
	var salt string
	var algorithm Algorithm
	err := a.DB.QueryRow(`
	SELECT 
		password_salt,
		algorithm
	FROM internal_logins WHERE user_id = $1`, uid).Scan(&salt, &algorithm)
	return salt, algorithm, err
}

func (a *AuthRepository) getInternalLoginUser(email string, hashedPassword string) (User, error) {
	var user User
	err := a.DB.QueryRow(`
		SELECT 
			u.id, 
			u.hash, 
			u.username, 
			u.email, 
			u.validated_at, 
			u.disabled, 
			u.created_at, 
			u.updated_at, 
			il.password_last_updated_at, 
			il.login_attempts, 
			il.last_login_attempt, 
			il.last_login_success,
			il.password_salt,
			il.algorithm,
			il.password
		FROM users u
		INNER JOIN internal_logins il ON u.id = il.user_id
		WHERE il.email = $1 AND il.password = $2
		`, email, hashedPassword).
		Scan(
			&user.ID,
			&user.Hash,
			&user.Username,
			&user.Email,
			&user.ValidatedAt,
			&user.Disabled,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.internalLoginData.PasswordLastUpdatedAt,
			&user.internalLoginData.LoginAttempts,
			&user.internalLoginData.LastLoginAttempt,
			&user.internalLoginData.LastLoginSuccess,
			&user.internalLoginData.PasswordSalt,
			&user.internalLoginData.Algorithm,
			&user.internalLoginData.Password,
		)
	return user, err
}

func (a *AuthRepository) createUser(username string, email string) (User, error) {
	var user User
	hash, err := uuid.NewV7()
	if err != nil {
		return user, err
	}
	err = a.DB.QueryRow(`
		INSERT INTO users (username, email, hash)
		VALUES ($1, $2, $3)
		RETURNING id, hash, username, email, validated_at, disabled, created_at, updated_at
		`, username, email, hash.String()).
		Scan(&user.ID, &user.Hash, &user.Username, &user.Email, &user.ValidatedAt, &user.Disabled, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (a *AuthRepository) createInternalLogin(uid int, hashedPassword string, salt string, algorithm uint16) (InternalLogin, error) {
	var il InternalLogin
	err := a.DB.QueryRow(`
		INSERT INTO internal_logins (user_id, email, password, password_salt, algorithm)
		VALUES ($1, (SELECT email FROM users WHERE id = $1), $2, $3, $4)
		RETURNING id, user_id, email, password, password_salt, algorithm, password_last_updated_at, login_attempts, last_login_attempt, last_login_success, created_at, updated_at
		`, uid, hashedPassword, salt, algorithm).
		Scan(&il.ID, &il.UserID, &il.Email, &il.Password, &il.PasswordSalt, &il.Algorithm, &il.PasswordLastUpdatedAt, &il.LoginAttempts, &il.LastLoginAttempt, &il.LastLoginSuccess, &il.CreatedAt, &il.UpdatedAt)
	return il, err
}

func (a *AuthRepository) increaseLoginAttempts(uid int) error {
	var attempts int
	err := a.DB.QueryRow(`
	UPDATE internal_logins
	SET login_attempts = login_attempts + 1
	WHERE user_id = $1
	RETURNING login_attempts
	`, uid).Scan(&attempts)
	if attempts >= 5 {
		log.Printf("Account locked for user %d", uid)
		_, err = a.DB.Exec(`
		UPDATE users
		SET disabled = true
		WHERE id = $1
		`, uid)
	}
	return err
}

func (a *AuthRepository) isAccountLocked(uid int) (bool, error) {
	var disabled bool
	err := a.DB.QueryRow("SELECT disabled FROM users WHERE id = $1", uid).Scan(&disabled)
	return disabled, err
}

func (a *AuthRepository) logUserSession(uid int, ip string, userAgent string) error {
	_, err := a.DB.Exec(`
	INSERT INTO user_sessions (user_id, login_ip, user_agent)
	VALUES ($1, $2, $3)
	`, uid, ip, userAgent)
	return err
}

func (a *AuthRepository) getUserByHash(hash string) (User, error) {
	var user User
	err := a.DB.QueryRow(`
		SELECT
			id,
			hash,
			username,
			email,
			validated_at,
			disabled,
			created_at,
			updated_at
		FROM users
		WHERE hash = $1
		`, hash).
		Scan(
			&user.ID,
			&user.Hash,
			&user.Username,
			&user.Email,
			&user.ValidatedAt,
			&user.Disabled,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	return user, err
}
