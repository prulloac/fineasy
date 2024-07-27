package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/prulloac/fineasy/pkg"
)

type InternalLogin struct {
	UserID                uint
	Password              string
	PasswordSalt          string
	Algorithm             pkg.Algorithm
	PasswordLastUpdatedAt time.Time
	LoginAttempts         int
	LastLoginAttempt      time.Time
	LastLoginSuccess      time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             sql.NullTime
}

func (i *InternalLogin) String() string {
	out, err := json.Marshal(i)
	if err != nil {
		return fmt.Sprintf("%+v", i.UserID)
	}
	return string(out)
}

func (a *AuthRepository) CreateInternalLogin(uid uint, hashedPassword string, salt string, algorithm pkg.Algorithm) (*InternalLogin, error) {
	var il InternalLogin
	err := a.Persistence.QueryRow(`
	INSERT INTO internal_logins 
	(user_id, password, password_salt, algorithm) VALUES ($1, $2, $3, $4) 
	RETURNING user_id, password, password_salt, algorithm, password_last_updated_at, login_attempts, last_login_attempt, last_login_success, created_at, updated_at, deleted_at
	`, uid, hashedPassword, salt, algorithm).
		Scan(&il.UserID, &il.Password, &il.PasswordSalt, &il.Algorithm, &il.PasswordLastUpdatedAt, &il.LoginAttempts, &il.LastLoginAttempt, &il.LastLoginSuccess, &il.CreatedAt, &il.UpdatedAt, &il.DeletedAt)
	return &il, err
}

func (a *AuthRepository) GetInternalLoginByUserID(uid uint) (InternalLogin, error) {
	var internalLogin InternalLogin
	err := a.Persistence.QueryRow(`
	SELECT 
		user_id, password, password_salt, algorithm, password_last_updated_at, login_attempts, last_login_attempt, last_login_success, created_at, updated_at, deleted_at
	FROM internal_logins 
	WHERE user_id = $1
	`, uid).Scan(&internalLogin.UserID, &internalLogin.Password, &internalLogin.PasswordSalt, &internalLogin.Algorithm, &internalLogin.PasswordLastUpdatedAt, &internalLogin.LoginAttempts, &internalLogin.LastLoginAttempt, &internalLogin.LastLoginSuccess, &internalLogin.CreatedAt, &internalLogin.UpdatedAt, &internalLogin.DeletedAt)
	return internalLogin, err
}

func (a *AuthRepository) IncreaseLoginAttempts(uid uint) error {
	var attempts int
	err := a.Persistence.QueryRow(`
	UPDATE internal_logins SET login_attempts = login_attempts + 1 WHERE user_id = $1 RETURNING login_attempts
	`, uid).Scan(&attempts)
	return err
}
