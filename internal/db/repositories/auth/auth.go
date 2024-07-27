package repositories

import (
	"github.com/prulloac/fineasy/internal/db/persistence"
	"github.com/prulloac/fineasy/pkg/logging"
)

type AuthRepository struct {
	Persistence *persistence.Persistence
	logger      *logging.Logger
}

func NewAuthRepository(persistence *persistence.Persistence) *AuthRepository {
	instance := &AuthRepository{}
	instance.Persistence = persistence
	instance.logger = logging.NewLoggerWithPrefix("[AuthRepository]")
	err := instance.CreateTables()
	if err != nil {
		panic(err)
	}
	return instance
}

func (a *AuthRepository) Close() {
	a.Persistence.Close()
}

func (a *AuthRepository) CreateTables() error {
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

func (a *AuthRepository) DropTables() error {
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
