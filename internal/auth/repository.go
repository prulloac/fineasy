package auth

import (
	"github.com/google/uuid"
	p "github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/pkg"
	"github.com/prulloac/fineasy/pkg/logging"
)

type Repository struct {
	Persistence *p.Persistence
	logger      *logging.Logger
}

func NewRepository(persistence *p.Persistence) *Repository {
	return &Repository{persistence, logging.NewLoggerWithPrefix("[Repository]")}
}

func (a *Repository) Close() {
	a.Persistence.Close()
}

func (a *Repository) CreateTables() error {
	r, err := a.Persistence.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		hash uuid NOT NULL,
		email VARCHAR(255) NOT NULL,
		validated_at TIMESTAMP,
		disabled BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS internal_logins (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL references users(id) ON DELETE CASCADE,
		password VARCHAR(255) NOT NULL,
		password_salt VARCHAR(255) NOT NULL,
		algorithm INTEGER NOT NULL,
		password_last_updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		login_attempts INT NOT NULL DEFAULT 0,
		last_login_attempt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		last_login_success TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP
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
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS external_login_tokens (
		id SERIAL PRIMARY KEY,
		external_login_id INT NOT NULL references external_logins(id) ON DELETE CASCADE,
		login_ip VARCHAR(255) NOT NULL,
		user_agent VARCHAR(255) NOT NULL,
		logged_in_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		token TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS user_sessions (
		id SERIAL PRIMARY KEY,
		session_token VARCHAR(255) NOT NULL,
		user_id INT NOT NULL references users(id) ON DELETE CASCADE,
		login_ip VARCHAR(255) NOT NULL,
		user_agent VARCHAR(255) NOT NULL,
		logged_in_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		logged_out_at TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE UNIQUE INDEX IF NOT EXISTS users_email_idx ON users (email);
	CREATE UNIQUE INDEX IF NOT EXISTS external_login_providers_name_idx ON external_login_providers (name, type, endpoint);
	CREATE UNIQUE INDEX IF NOT EXISTS external_logins_user_id_provider_id_idx ON external_logins (user_id, provider_id);
	`)
	if err != nil {
		a.logger.Fatalf("Error creating tables: %s", err)
		return err
	}
	a.logger.Printf("Tables created: %v", r)
	return nil
}

func (a *Repository) DropTables() error {
	_, err := a.Persistence.Exec(`
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

func (a *Repository) getUserIDByEmail(email string) (uint, error) {
	var uid uint
	err := a.Persistence.QueryRow(`
	SELECT id FROM users WHERE email = $1
	`, email).Scan(&uid)
	if err != nil {
		a.logger.Printf("Error getting user by email: %s for email: %s", err, email)
		return 0, err
	}
	return uid, nil
}

func (a *Repository) getSaltAndAlgorithmByUserID(uid uint) (string, pkg.Algorithm, error) {
	// salt, algorithm,
	var sa struct {
		PasswordSalt string
		Algorithm    pkg.Algorithm
	}
	err := a.Persistence.QueryRow(`
	SELECT password_salt, algorithm FROM internal_logins WHERE user_id = $1
	`, uid).
		Scan(&sa.PasswordSalt, &sa.Algorithm)
	return sa.PasswordSalt, sa.Algorithm, err
}

func (a *Repository) getInternalLoginUserByEmailAndPassword(email string, hashedPassword string) (*User, error) {
	var user User
	err := a.Persistence.QueryRow(`
	SELECT 
		users.id, users.hash, users.email, users.validated_at, users.disabled, users.created_at, users.updated_at, users.deleted_at
	FROM users
	JOIN internal_logins ON users.id = internal_logins.user_id
	WHERE users.email = $1 AND internal_logins.password = $2
	`, email, hashedPassword).
		Scan(&user.ID, &user.Hash, &user.Email, &user.ValidatedAt, &user.Disabled, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	return &user, err
}

func (a *Repository) createUser(email string) (*User, error) {
	hash, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	var user User
	err = a.Persistence.QueryRow(`
	INSERT INTO users 
	(email, hash) VALUES ($1, $2)
	RETURNING id, hash, email, validated_at, disabled, created_at, updated_at, deleted_at
	`, email, hash.String()).
		Scan(&user.ID, &user.Hash, &user.Email, &user.ValidatedAt, &user.Disabled, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	return &user, err
}

func (a *Repository) createInternalLogin(uid uint, hashedPassword string, salt string, algorithm pkg.Algorithm) (*InternalLogin, error) {
	var il InternalLogin
	err := a.Persistence.QueryRow(`
	INSERT INTO internal_logins 
	(user_id, password, password_salt, algorithm) VALUES ($1, $2, $3, $4) 
	RETURNING id, user_id, password, password_salt, algorithm, password_last_updated_at, login_attempts, last_login_attempt, last_login_success, created_at, updated_at, deleted_at
	`, uid, hashedPassword, salt, algorithm).
		Scan(&il.ID, &il.UserID, &il.Password, &il.PasswordSalt, &il.Algorithm, &il.PasswordLastUpdatedAt, &il.LoginAttempts, &il.LastLoginAttempt, &il.LastLoginSuccess, &il.CreatedAt, &il.UpdatedAt, &il.DeletedAt)
	return &il, err
}

func (a *Repository) increaseLoginAttempts(uid uint) error {
	var attempts int
	err := a.Persistence.QueryRow(`
	UPDATE internal_logins SET login_attempts = login_attempts + 1 WHERE user_id = $1 RETURNING login_attempts
	`, uid).Scan(&attempts)
	return err
}

func (a *Repository) isAccountLocked(uid uint) (bool, error) {
	var disabled bool
	err := a.Persistence.QueryRow(`
	SELECT disabled FROM users WHERE id = $1
	`, uid).Scan(&disabled)
	return disabled, err
}

func (a *Repository) logUserSession(uid uint, ip string, userAgent string) (*UserSession, error) {
	token, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	var session UserSession
	err = a.Persistence.QueryRow(`
	INSERT INTO user_sessions
	(user_id, login_ip, user_agent, session_token) VALUES ($1, $2, $3, $4)
	RETURNING id, user_id, login_ip, user_agent, session_token, logged_in_at, logged_out_at, created_at, updated_at
	`, uid, ip, userAgent, token.String()).Scan(&session.ID, &session.UserID, &session.LoginIP, &session.UserAgent, &session.SessionToken, &session.LoggedInAt, &session.LoggedOutAt, &session.CreatedAt, &session.UpdatedAt)
	return &session, err
}
