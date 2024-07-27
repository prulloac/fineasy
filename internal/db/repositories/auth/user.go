package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/prulloac/fineasy/pkg"
)

type User struct {
	pkg.Model
	Hash              string
	Email             string
	ValidatedAt       sql.NullTime
	Disabled          bool
	InternalLoginData InternalLogin
	ExternalLoginData ExternalLogin
}

func (u *User) String() string {
	out, err := json.Marshal(u)
	if err != nil {
		return fmt.Sprintf("%+v", u.Email)
	}
	return string(out)
}

func (a *AuthRepository) CreateUser(email string) (*User, error) {
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

func (a *AuthRepository) GetUserIDByEmail(email string) (uint, error) {
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

func (a *AuthRepository) IsAccountLocked(uid uint) (bool, error) {
	var disabled bool
	err := a.Persistence.QueryRow(`
	SELECT disabled FROM users WHERE id = $1
	`, uid).Scan(&disabled)
	return disabled, err
}

func (a *AuthRepository) GetUserByEmailAndPassword(mail, pwd string) (*User, error) {
	var user User
	err := a.Persistence.Debug().QueryRow(`
	SELECT u.id, u.hash, u.email, u.validated_at, u.disabled, u.created_at, u.updated_at, u.deleted_at,
		il.user_id, il.password, il.password_salt, il.algorithm, il.created_at, il.updated_at, il.deleted_at
	FROM users u
	JOIN internal_logins il ON u.id = il.user_id
	WHERE u.email = $1
	`, mail).Scan(&user.ID, &user.Hash, &user.Email, &user.ValidatedAt, &user.Disabled, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		&user.InternalLoginData.UserID, &user.InternalLoginData.Password, &user.InternalLoginData.PasswordSalt, &user.InternalLoginData.Algorithm, &user.InternalLoginData.CreatedAt, &user.InternalLoginData.UpdatedAt, &user.InternalLoginData.DeletedAt)
	return &user, err
}
